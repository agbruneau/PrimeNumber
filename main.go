/*
 * Fichier: main.go
 * Auteur: [Votre Nom/Organisation]
 * Date: 20 juin 2025
 *
 * Description:
 * Ce programme est une vérification empirique et une implémentation optimisée
 * pour trouver des nombres premiers 'n' qui satisfont au théorème prouvé par
 * les mathématiciens Ben Green et Mehtaab Sawhney.
 *
 * Le théorème stipule qu'il existe une infinité de nombres premiers de la forme:
 * n = p^2 + 4*q^2
 * où 'p' et 'q' sont eux-mêmes des nombres premiers.
 *
 * Architecture de la solution:
 * - Utilisation d'un crible d'Eratosthène pour la génération efficace des nombres premiers initiaux.
 * - Implémentation d'un pool de workers (Worker Pool) avec des goroutines pour paralléliser
 * la recherche et tirer parti des processeurs multi-cœurs.
 * - Utilisation de canaux (channels) pour la distribution des tâches et la collecte des résultats
 * de manière concurrente et sécurisée.
 */
package main

import (
	"flag"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

// Job représente une tâche à effectuer par un worker: une paire (p, q) à tester.
type Job struct {
	p int
	q int
}

// Result représente un résultat positif trouvé par un worker.
type Result struct {
	p int
	q int
	n int
}

// sieveOfEratosthenes génère tous les nombres premiers jusqu'à une limite donnée.
// C'est une méthode beaucoup plus efficace que des tests de primalité individuels.
func sieveOfEratosthenes(limit int) []int {
	// Initialise un tableau de booléens pour marquer les nombres.
	// `primes[i]` sera `true` si `i` n'est pas premier.
	primesMarker := make([]bool, limit+1)
	primesMarker[0], primesMarker[1] = true, true // 0 et 1 ne sont pas premiers.

	// Algorithme du crible.
	for p := 2; p*p <= limit; p++ {
		if !primesMarker[p] { // Si p est premier...
			for i := p * p; i <= limit; i += p {
				primesMarker[i] = true // ...marquer tous ses multiples comme non premiers.
			}
		}
	}

	// Collectionner les nombres premiers.
	// Pre-allocate primes slice with an estimated capacity to reduce reallocations.
	// Prime Number Theorem: pi(x) ~ x / ln(x)
	var estimatedPrimes int
	if limit < 2 { // Avoid Log of numbers < 1, and handle small limits
		estimatedPrimes = 0
	} else if limit < 20 { // For very small limits, ln(limit) can be too small or estimation poor.
		estimatedPrimes = limit // Overestimate slightly or use fixed small capacity
	} else {
		estimatedPrimes = int(float64(limit) / math.Log(float64(limit)))
	}
	// Ensure a minimum capacity, e.g. for limit=2, estimate is ~1, actual is 1. For limit=10, estimate ~4, actual 4.
	// Give some extra capacity factor, e.g., 1.2, because estimate is an approximation.
	primes := make([]int, 0, int(float64(estimatedPrimes)*1.2)+10) // +10 as a small buffer

	for p := 2; p <= limit; p++ {
		if !primesMarker[p] {
			primes = append(primes, p)
		}
	}
	if len(primes) == 0 {
		return nil
	}
	return primes
}

// isPrime vérifie si un grand nombre est premier.
// Nécessaire pour les résultats 'n' qui peuvent dépasser la limite du crible.
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i = i + 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// power calculates (base^exp) % mod efficiently.
// Necessary for Miller-Rabin.
func power(base, exp, mod int) int {
	res := 1
	base %= mod
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return res
}

// isPrimeMillerRabin implements the Miller-Rabin primality test.
// k is the number of rounds for testing. Higher k means more accuracy.
// For a deterministic version for numbers up to 2^64, specific bases can be used.
// Here, we'll use k random bases for simplicity, good for typical int sizes.
// Returns true if n is likely prime, false if composite.
func isPrimeMillerRabin(n int, k int) bool {
	if n <= 1 || n == 4 {
		return false
	}
	if n <= 3 { // 2 and 3
		return true
	}
	if n%2 == 0 {
		return false
	}

	// Write n-1 as 2^s * d
	d := n - 1
	s := 0
	for d%2 == 0 {
		d /= 2
		s++
	}

	// Witness loop
	// Using math/rand for simplicity. For cryptographic purposes, crypto/rand is needed.
	// Seed is managed globally or passed around. For this use case,
	// time-based seeding in main or once globally is sufficient.
	// Since we don't have direct access to main's seeding here,
	// this might produce same random numbers if called very rapidly in parallel
	// without external seeding. However, worker calls are somewhat spread out.
	// For now, let's assume seeding is handled externally or this is acceptable.
	// A more robust way would be to pass a *rand.Rand source.

	for i := 0; i < k; i++ {
		// Pick a random 'a' in [2, n-2]
		// To avoid issues with rand.Intn(0) for n=2 or n=3 (already handled),
		// and to ensure a is in [2, n-2].
		// rand.Intn(max-min+1) + min
		a := 2 + int(time.Now().UnixNano())%(n-3) // Not cryptographically secure random.
		// A simpler way for non-crypto rand: a := rand.Intn(n-3) + 2

		x := power(a, d, n)

		if x == 1 || x == n-1 {
			continue
		}

		witness := true
		for r := 1; r < s; r++ {
			x = power(x, 2, n)
			if x == n-1 {
				witness = false
				break
			}
		}
		if witness {
			return false // n is composite
		}
	}
	return true // n is probably prime
}

// worker est une fonction qui s'exécute dans une goroutine.
// Elle reçoit des tâches (Jobs) depuis un canal, les traite,
// et envoie les résultats positifs dans un autre canal.
// Le paramètre 'id' est ignoré avec '_' pour résoudre l'alerte du linter.
// worker now takes primeTestAlgorithm and millerRabinK to decide which primality test to use.
func worker(_ int, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result, primeTestAlgorithm string, millerRabinK int) {
	defer wg.Done() // Signale que ce worker a terminé lorsque la fonction retourne.

	for job := range jobs { // Itère sur le canal de tâches jusqu'à sa fermeture.
		p, q := job.p, job.q
		n := (p * p) + 4*(q*q)

		var currentIsPrime bool
		if primeTestAlgorithm == "miller" {
			currentIsPrime = isPrimeMillerRabin(n, millerRabinK)
		} else { // Default or "trial"
			currentIsPrime = isPrime(n)
		}

		if currentIsPrime {
			results <- Result{p: p, q: q, n: n}
		}
	}
}

func main() {
	startTime := time.Now()

	// --- Configuration ---
	// Définition du flag pour searchLimit
	// Le premier argument est le nom du flag.
	// Le deuxième est la valeur par défaut.
	// Le troisième est la description du flag (utilisée par -help).
	searchLimitPtr := flag.Int("limit", 1000, "Limite supérieure pour la recherche des nombres premiers p et q.")
	primeTestPtr := flag.String("primetest", "trial", "Algorithme de test de primalité à utiliser: 'trial' ou 'miller'.")
	millerRabinIterationsPtr := flag.Int("k", 5, "Nombre d'itérations pour Miller-Rabin (si utilisé).")

	flag.Parse() // Analyse les arguments de la ligne de commande.

	searchLimit := *searchLimitPtr // Déréférence le pointeur pour obtenir la valeur.
	primeTestAlgorithm := *primeTestPtr // Value moved up, already applied
	millerRabinK := *millerRabinIterationsPtr // Value moved up, already applied

	// Utilisation de tous les cœurs de processeur disponibles pour les workers.
	numWorkers := runtime.NumCPU()

	fmt.Printf("Initialisation avec searchLimit=%d, numWorkers=%d, primeTest='%s'\n", searchLimit, numWorkers, primeTestAlgorithm)
	if primeTestAlgorithm == "miller" {
		fmt.Printf("Miller-Rabin itérations k=%d\n", millerRabinK)
	}
	fmt.Println("-------------------------------------------------------------------")

	// --- Étape 1: Génération optimisée des nombres premiers ---
	fmt.Println("Génération des nombres premiers avec le crible d'Eratosthène...")
	primes := sieveOfEratosthenes(searchLimit)
	fmt.Printf("%d nombres premiers trouvés jusqu'à %d.\n\n", len(primes), searchLimit)

	// --- Étape 2: Mise en place du Pool de Workers et des canaux ---
	jobs := make(chan Job, len(primes))
	results := make(chan Result, 100) // Canal avec buffer pour les résultats.
	var wg sync.WaitGroup

	// Démarrage des workers.
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		// Pass primeTestAlgorithm and millerRabinK to each worker
		go worker(w, &wg, jobs, results, primeTestAlgorithm, millerRabinK)
	}

	// --- Étape 3: Distribution des tâches ---
	// Une goroutine distincte est utilisée pour envoyer les tâches afin de ne pas bloquer
	// la collecte des résultats, qui se fait en parallèle.
	go func() {
		for _, p := range primes {
			for _, q := range primes {
				jobs <- Job{p: p, q: q}
			}
		}
		close(jobs) // Ferme le canal, signale aux workers qu'il n'y a plus de tâches.
	}()

	// --- Étape 4: Collecte des résultats ---
	// Une goroutine pour fermer le canal de résultats une fois que tous les workers ont terminé.
	go func() {
		wg.Wait() // Attend la fin de tous les workers.
		close(results)
	}()

	// Affichage des résultats au fur et à mesure de leur arrivée.
	fmt.Printf("%-10s | %-10s | %-20s | %-s\n", "p", "q", "n = p^2 + 4q^2", "Vérification")
	count := 0
	for res := range results {
		count++
		fmt.Printf("%-10d | %-10d | %-20d | %s\n", res.p, res.q, res.n, "Trouvé!")
	}

	// --- Finalisation ---
	duration := time.Since(startTime)
	fmt.Println("-------------------------------------------------------------------")
	fmt.Printf("Recherche terminée. %d nombres premiers spéciaux trouvés.\n", count)
	fmt.Printf("\nDurée totale de l'exécution: %s\n", duration)
}
