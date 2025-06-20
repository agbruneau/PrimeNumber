/*
 * Fichier: main_test.go
 * Auteur: [Votre Nom/Organisation]
 * Date: 20 juin 2025
 *
 * Description:
 * Ce fichier contient les tests unitaires pour le programme de vérification
 * du théorème sur les nombres premiers. Il valide le crible d'Eratosthène
 * et les différentes fonctions de test de primalité.
 */
package main

import (
	"reflect"
	"testing"
)

// TestSieveOfEratosthenes valide la génération des nombres premiers.
func TestSieveOfEratosthenes(t *testing.T) {
	testCases := []struct {
		name     string
		limit    int
		expected []int
	}{
		{"Limite de 30", 30, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29}},
		{"Limite de 10", 10, []int{2, 3, 5, 7}},
		{"Limite de 2", 2, []int{2}},
		{"Limite de 1", 1, nil},
		{"Limite de 0", 0, nil},
		{"Limite négative", -10, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := sieveOfEratosthenes(tc.limit)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("sieveOfEratosthenes(%d) = %v, attendu %v", tc.limit, result, tc.expected)
			}
		})
	}
}

// TestIsNPrimeAccordingToGreenSawhneyContext valide le test de primalité par division utilisé dans le contexte de Green-Sawhney.
func TestIsNPrimeAccordingToGreenSawhneyContext(t *testing.T) {
	testCases := []struct {
		name     string
		n        int64
		expected bool
	}{
		{"Nombre premier 2", 2, true},
		{"Nombre premier 7", 7, true},
		{"Nombre premier 97", 97, true},
		{"Nombre composé 1", 1, false},
		{"Nombre composé 4", 4, false},
		{"Nombre composé 100", 100, false},
		{"Grand nombre premier", 7919, true},
		{"Grand nombre composé", 7921, false}, // 89*89
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := isNPrimeAccordingToGreenSawhneyContext(tc.n); result != tc.expected {
				t.Errorf("isNPrimeAccordingToGreenSawhneyContext(%d) = %v, attendu %v", tc.n, result, tc.expected)
			}
		})
	}
}

// TestIsPrimeMillerRabin64 valide le test de primalité de Miller-Rabin.
func TestIsPrimeMillerRabin64(t *testing.T) {
	testCases := []struct {
		name     string
		n        int64
		expected bool
	}{
		{"Nombre premier 2", 2, true},
		{"Nombre premier 7", 7, true},
		{"Nombre premier 97", 97, true},
		{"Nombre composé 1", 1, false},
		{"Nombre composé 4", 4, false},
		{"Nombre composé 100", 100, false},
		{"Grand nombre premier", 7919, true},
		{"Grand nombre composé", 7921, false}, // 89*89
		// Test avec un grand nombre premier pour lequel Miller-Rabin est plus efficace
		{"Très grand nombre premier", 2147483647, true}, // Nombre premier de Mersenne
		{"Très grand nombre composé", 2147483649, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := isPrimeMillerRabin64(tc.n); result != tc.expected {
				t.Errorf("isPrimeMillerRabin64(%d) = %v, attendu %v", tc.n, result, tc.expected)
			}
		})
	}
}

// TestPower64 valide la fonction d'exponentiation modulaire.
func TestPower64(t *testing.T) {
	testCases := []struct {
		name           string
		base, exp, mod int64
		expected       int64
	}{
		{"Simple", 2, 3, 10, 8},
		{"Modulo s'applique", 5, 2, 20, 5},
		// Correction de la valeur attendue de 467 à 43.
		{"Grand nombres", 123, 45, 1000, 43},
		{"Exp 0", 123, 0, 1000, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if result := power64(tc.base, tc.exp, tc.mod); result != tc.expected {
				t.Errorf("power64(%d, %d, %d) = %d, attendu %d", tc.base, tc.exp, tc.mod, result, tc.expected)
			}
		})
	}
}
