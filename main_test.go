/*
 * Fichier: main_test.go
 * Auteur: [Votre Nom/Organisation]
 * Date: 20 juin 2025
 *
 * Description:
 * Ce fichier contient les tests unitaires et les benchmarks pour le programme
 * principal qui vérifie le théorème de Green-Sawhney. Les tests se concentrent
 * sur la validation de la logique métier principale: la génération de nombres
 * premiers et la vérification de la primalité.
 */
package main

import (
	"reflect"
	"testing"
)

// TestSieveOfEratosthenes valide la correction de l'algorithme du crible d'Eratosthène.
func TestSieveOfEratosthenes(t *testing.T) {
	// Définition de cas de test avec une approche de "table-driven tests".
	testCases := []struct {
		name     string // Nom du cas de test.
		limit    int    // Limite pour la génération des nombres premiers.
		expected []int  // Le résultat attendu.
	}{
		{
			name:     "Limite de 30",
			limit:    30,
			expected: []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29},
		},
		{
			name:     "Limite de 10",
			limit:    10,
			expected: []int{2, 3, 5, 7},
		},
		{
			name:     "Limite de 2",
			limit:    2,
			expected: []int{2},
		},
		{
			name:     "Limite inférieure (1)",
			limit:    1,
			expected: nil, // Un slice nil est retourné pour une limite < 2.
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Exécution de la fonction à tester.
			result := sieveOfEratosthenes(tc.limit)
			// reflect.DeepEqual est utilisé pour comparer des slices (tableaux dynamiques).
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Pour la limite %d, résultat obtenu %#v, attendu %#v", tc.limit, result, tc.expected)
			}
		})
	}
}

// TestIsPrime valide la correction de la fonction de vérification de primalité.
func TestIsPrime(t *testing.T) {
	testCases := []struct {
		name     string // Nom du cas de test.
		number   int    // Nombre à tester.
		expected bool   // Résultat de primalité attendu.
	}{
		{"Nombre premier (2)", 2, true},
		{"Nombre premier (3)", 3, true},
		{"Nombre premier (7)", 7, true},
		{"Nombre premier (41)", 41, true},
		{"Grand nombre premier", 7919, true},
		{"Nombre non premier (0)", 0, false},
		{"Nombre non premier (1)", 1, false},
		{"Nombre non premier (4)", 4, false},
		{"Nombre non premier (25)", 25, false},
		{"Grand nombre non premier", 8000, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := isPrime(tc.number); result != tc.expected {
				t.Errorf("Pour le nombre %d, primalité obtenue %t, attendue %t", tc.number, result, tc.expected)
			}
		})
	}
}

// BenchmarkSieveOfEratosthenes mesure la performance de la génération de nombres premiers.
func BenchmarkSieveOfEratosthenes(b *testing.B) {
	// Le framework de benchmark exécute la boucle b.N fois pour obtenir une mesure stable.
	for i := 0; i < b.N; i++ {
		sieveOfEratosthenes(10000) // Utilise une limite raisonnable pour le benchmark.
	}
}

// BenchmarkIsPrime mesure la performance du test de primalité pour un grand nombre.
func BenchmarkIsPrime(b *testing.B) {
	// Test avec un grand nombre premier.
	largePrime := 999999999989
	for i := 0; i < b.N; i++ {
		isPrime(largePrime)
	}
}
