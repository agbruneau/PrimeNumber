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
 * - Utilisation de types int64 et du paquet math/big pour garantir l'exactitude avec de grands nombres.
 */
package main

import (
	"flag"
	"fmt"
	"math"
	"math/big"
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
// Le type de 'n' est int64 pour éviter les débordements (overflows).
type Result struct {
	p int
	q int
	n int64
}

// sieveOfEratosthenes génère tous les nombres premiers jusqu'à une limite donnée.
// C'est une méthode beaucoup plus efficace que des tests de primalité individuels.
func sieveOfEratosthenes(limit int) []int {
	// Ajout d'une validation pour gérer les cas limites (négatifs, 0, 1)
	// et prévenir les erreurs "index out of range".
	if limit < 2 {
		return nil
	}

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
	// Pré-allouer la slice de nombres premiers avec une capacité estimée pour réduire les réallocations.
	// Théorème des nombres premiers: pi(x) ~ x / ln(x)
	var estimatedPrimes int
	if limit > 1 {
		estimatedPrimes = int(float64(limit) / math.Log(float64(limit)))
	}
	primes := make([]int, 0, int(float64(estimatedPrimes)*1.2)+10)

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

// isNPrimeAccordingToGreenSawhneyContext vérifie si un grand nombre est premier par division successive.
// Ce nom reflète son utilisation dans le contexte de la vérification des nombres 'n' issus
// de la formule p^2 + 4q^2 du théorème de Green-Sawhney.
// L'algorithme sous-jacent reste la division par essais.
// Nécessaire pour les résultats 'n' qui peuvent dépasser la limite du crible.
// Utilise int64 pour la robustesse.
func isNPrimeAccordingToGreenSawhneyContext(n int64) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	// On vérifie les diviseurs de la forme 6k ± 1 jusqu'à sqrt(n).
	limit := int64(math.Sqrt(float64(n)))
	for i := int64(5); i <= limit; i = i + 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// power64 calcule (base^exp) % mod de manière sûre avec math/big pour éviter les débordements.
func power64(base, exp, mod int64) int64 {
	bBase := big.NewInt(base)
	bExp := big.NewInt(exp)
	bMod := big.NewInt(mod)

	// Le paquet math/big gère les grands nombres de manière sûre.
	res := new(big.Int)
	res.Exp(bBase, bExp, bMod)

	return res.Int64()
}

// isPrimeMillerRabin64 implémente le test de primalité de Miller-Rabin.
// Cette version est déterministe pour tous les nombres de type int64.
// Elle utilise un ensemble de bases prédéfinies qui garantissent l'exactitude.
func isPrimeMillerRabin64(n int64) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	// Écrire n-1 comme 2^s * d
	d := n - 1
	s := 0
	for d%2 == 0 {
		d /= 2
		s++
	}

	// Bases de test qui rendent l'algorithme déterministe pour n < 2^64.
	bases := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	// Pour n < 3,317,044,064,279,371, les 12 premières bases suffisent.

	for _, a := range bases {
		if a >= n-1 {
			break
		}
		x := power64(a, d, n)

		if x == 1 || x == n-1 {
			continue
		}

		isWitness := true
		for r := 1; r < s; r++ {
			x = power64(x, 2, n)
			if x == n-1 {
				isWitness = false
				break
			}
		}

		if isWitness {
			return false // n est composé.
		}
	}

	return true // n est probablement (ici, certainement) premier.
}

// worker est une fonction qui s'exécute dans une goroutine.
// Elle reçoit des tâches (Jobs) depuis un canal, les traite,
// et envoie les résultats positifs dans un autre canal.
func worker(wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result, primeTestAlgorithm string) {
	defer wg.Done()

	for job := range jobs {
		p, q := int64(job.p), int64(job.q)
		n := (p * p) + 4*(q*q)

		var isNPrime bool
		if primeTestAlgorithm == "miller" {
			isNPrime = isPrimeMillerRabin64(n)
		} else { // Par défaut: "trial"
			isNPrime = isNPrimeAccordingToGreenSawhneyContext(n)
		}

		if isNPrime {
			results <- Result{p: job.p, q: job.q, n: n}
		}
	}
}

func main() {
	startTime := time.Now()

	// --- Configuration ---
	searchLimitPtr := flag.Int("limit", 1000, "Limite supérieure pour la recherche des nombres premiers p et q.")
	primeTestPtr := flag.String("primetest", "miller", "Algorithme de test de primalité: 'trial' ou 'miller' (défaut).")
	flag.Parse()

	searchLimit := *searchLimitPtr
	primeTestAlgorithm := *primeTestPtr

	numWorkers := runtime.NumCPU()

	fmt.Printf("Initialisation avec searchLimit=%d, numWorkers=%d, primeTest='%s'\n", searchLimit, numWorkers, primeTestAlgorithm)
	fmt.Println("-------------------------------------------------------------------")

	// --- Étape 1: Génération optimisée des nombres premiers ---
	fmt.Println("Génération des nombres premiers avec le crible d'Eratosthène...")
	primes := sieveOfEratosthenes(searchLimit)
	if primes == nil {
		fmt.Println("Aucun nombre premier trouvé dans la limite spécifiée.")
		return
	}
	fmt.Printf("%d nombres premiers trouvés jusqu'à %d.\n\n", len(primes), searchLimit)

	// --- Étape 2: Mise en place du Pool de Workers et des canaux ---
	jobs := make(chan Job, len(primes))
	results := make(chan Result, 100)
	var wg sync.WaitGroup

	// Démarrage des workers.
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(&wg, jobs, results, primeTestAlgorithm)
	}

	// --- Étape 3: Distribution des tâches ---
	go func() {
		for _, p := range primes {
			for _, q := range primes {
				jobs <- Job{p: p, q: q}
			}
		}
		close(jobs) // Ferme le canal, signale aux workers qu'il n'y a plus de tâches.
	}()

	// --- Étape 4: Collecte des résultats ---
	go func() {
		wg.Wait() // Attend la fin de tous les workers.
		close(results)
	}()

	fmt.Printf("%-10s | %-10s | %-25s | %-s\n", "p", "q", "n = p^2 + 4q^2", "Vérification")
	count := 0
	for res := range results {
		count++
		fmt.Printf("%-10d | %-10d | %-25d | %s\n", res.p, res.q, res.n, "Trouvé!")
	}

	// --- Finalisation ---
	duration := time.Since(startTime)
	fmt.Println("-------------------------------------------------------------------")
	fmt.Printf("Recherche terminée. %d nombres premiers spéciaux trouvés.\n", count)
	fmt.Printf("\nDurée totale de l'exécution: %s\n", duration)
}
