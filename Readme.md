# Vérification du Théorème de Green-Sawhney sur les Nombres Premiers

Ce programme Go est une implémentation conçue pour rechercher empiriquement des nombres premiers `n` qui satisfont à la formule `n = p^2 + 4*q^2`, où `p` et `q` sont eux-mêmes des nombres premiers. Ceci est lié à un théorème prouvé par les mathématiciens Ben Green et Mehtaab Sawhney, qui affirme l'existence d'une infinité de tels nombres premiers.

## Fonctionnalités et Optimisations

*   **Génération Efficace de Nombres Premiers**: Utilise le **Crible d'Eratosthène** pour générer rapidement la liste initiale des nombres premiers `p` et `q` jusqu'à une limite spécifiée.
*   **Traitement Parallèle**: Met en œuvre un **pool de workers (Worker Pool)** utilisant des goroutines Go pour paralléliser la vérification des paires `(p, q)`. Cela permet de tirer parti des processeurs multi-cœurs et d'accélérer considérablement la recherche.
*   **Communication Concurrente Sécurisée**: Utilise des canaux (channels) Go pour distribuer les tâches aux workers et collecter les résultats de manière sûre en concurrence.
*   **Test de Primalité Optimisé**: La fonction `isPrime` utilisée pour vérifier la primalité des grands nombres `n` (résultats de `p^2 + 4*q^2`) est optimisée pour ignorer les multiples de 2 et 3, et ne vérifier que les diviseurs de la forme `6k ± 1`.

## Prérequis

*   Go (version 1.21 ou plus récente recommandée)

## Compilation et Exécution

1.  **Cloner le dépôt (si applicable) ou télécharger les fichiers `main.go`, `main_test.go`, et `go.mod` dans un répertoire.**

2.  **Ouvrir un terminal et naviguer vers le répertoire du projet.**

3.  **Compiler le programme :**
    ```bash
    go build
    ```
    Cela créera un exécutable nommé `PrimeNumber` (ou `PrimeNumber.exe` sous Windows).

4.  **Exécuter le programme :**
    *   Pour utiliser la limite de recherche par défaut (actuellement codée en dur à `1000` pour `p` et `q`):
        ```bash
        ./PrimeNumber
        ```
    *   (Après l'implémentation du flag à l'étape suivante du plan) Pour spécifier une limite de recherche différente (par exemple, 500) :
        ```bash
        ./PrimeNumber -limit=500
        ```

    Le programme affichera les nombres premiers `p` et `q` trouvés, ainsi que le nombre `n` résultant qui est également premier. Il indiquera également le nombre total de ces nombres premiers spéciaux trouvés et la durée totale de l'exécution.

## Exécution des Tests

Pour exécuter les tests unitaires et les benchmarks :

```bash
go test ./...
```

Pour exécuter les benchmarks avec plus de détails (par exemple, l'allocation mémoire) :

```bash
go test -bench=. -benchmem ./...
```

## Structure du Code

*   `main.go`: Contient la logique principale du programme, y compris le crible d'Eratosthène, la fonction de test de primalité, la gestion du pool de workers, et la fonction `main`.
*   `main_test.go`: Contient les tests unitaires pour les fonctions `sieveOfEratosthenes` et `isPrime`, ainsi que des benchmarks de performance.
*   `go.mod`: Définit le module Go et ses dépendances (aucune dépendance externe pour le moment).
*   `Readme.md`: Ce fichier.

## Auteur

Ce programme a été adapté et optimisé. L'auteur original du code de base n'est pas spécifié, mais les améliorations et la structuration actuelles sont le résultat de ce projet.
