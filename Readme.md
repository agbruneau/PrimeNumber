package main

import (
	"fmt"
	"time"
)

// isPrime vérifie si un nombre est premier.
// Utilise une approche optimisée de division par essais.
// Un nombre est premier s'il n'est divisible que par 1 et par lui-même.
// On vérifie la divisibilité de 2 jusqu'à la racine carrée du nombre.
func isPrime(n int) bool {
	// 0 et 1 ne sont pas des nombres premiers.
	if n <= 1 {
		return false
	}
	// 2 et 3 sont des nombres premiers.
	if n <= 3 {
		return true
	}
	// Si le nombre est divisible par 2 ou 3, il n'est pas premier.
	// C'est une optimisation pour éliminer rapidement de nombreux composites.
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	// Tous les nombres premiers (sauf 2 et 3) sont de la forme 6k ± 1.
	// Nous pouvons donc vérifier les diviseurs en sautant de 6 en 6.
	for i := 5; i*i <= n; i = i + 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// findSpecialPrimes recherche les nombres premiers de la forme p^2 + 4q^2.
// 'limit' définit la borne supérieure pour les nombres premiers p et q.
func findSpecialPrimes(limit int) {
	fmt.Printf("Recherche des nombres premiers de la forme n = p^2 + 4q^2 jusqu'à p, q <= %d\n", limit)
	fmt.Println("-------------------------------------------------------------------")

	// Étape 1: Générer une liste de nombres premiers jusqu'à la limite.
	// Cela évite de recalculer la primalité de p et q à chaque itération.
	primes := []int{}
	for i := 2; i <= limit; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}

	fmt.Printf("Trouvé %d nombres premiers jusqu'à %d. Début de la recherche...\n\n", len(primes), limit)
	fmt.Printf("%-10s | %-10s | %-20s | %-s\n", "p", "q", "n = p^2 + 4q^2", "Vérification")

	count := 0
	// Étape 2: Itérer sur toutes les paires possibles de nombres premiers (p, q).
	for _, p := range primes {
		for _, q := range primes {
			// Calculer la valeur de n selon la formule.
			n := (p * p) + 4*(q*q)

			// Étape 3: Vérifier si le résultat n est également un nombre premier.
			if isPrime(n) {
				count++
				fmt.Printf("%-10d | %-10d | %-20d | %s\n", p, q, n, "Trouvé!")
			}
		}
	}

	fmt.Println("-------------------------------------------------------------------")
	fmt.Printf("Recherche terminée. %d nombres premiers spéciaux trouvés.\n", count)
}

func main() {
	// Démarrer un chronomètre pour mesurer la durée d'exécution.
	startTime := time.Now()

	// Définir la limite pour la recherche des nombres premiers p et q.
	// Attention: une limite élevée augmentera considérablement le temps de calcul (complexité O(N^2)).
	// Une limite de 500 est raisonnable pour une exécution rapide.
	searchLimit := 500

	findSpecialPrimes(searchLimit)

	// Calculer et afficher la durée totale de l'exécution.
	duration := time.Since(startTime)
	fmt.Printf("\nDurée totale de l'exécution: %s\n", duration)
}
