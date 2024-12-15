#notes #books #compititive #programming #graph #trees

# Graph / Trees Algorithms
Remember tree is just a graph with only one unique path between two nodes

```
Travesals
1. BFS
2. DFS

Shortest Path
1. Djikstra's
2. Bellman Ford
3. Floydd Warshall

DataStructure to Union and Find
1. Union Find / Disjoint Set

Minimum/maximum Spanning Tree
1.Kruskal's Algo
2.Prim's Algo

Cycle in a Directed Graph
1. Flyod'S Algorithm (fast / slow pointer)

Strongly Connected Graphs
1. Kosaraju's Algorithm
2. 2SAT Problem/Algorithm
```

## 2SAT Problem

problem of assigning Boolean values to variables to satisfy a formulas
for eg, assign values to xi so the the result in true

```
(a & ~b) OR (~a & b) OR (~a & ~b) OR (a & ~c)
```

basically we build a connected graph out of the expression so for a V b
there is an edge from  ~a -> b / ~b -> a

which means if a is false, b should be true. B is false a should be true

which helps us determine that if a can reach ~a  and vice versa then there exists a path where both can be true (which is not possible ~a is negate of a).
and hence the formula doesnt not have any solution

essentially this would happen if there are part of the same connected component

https://cp-algorithms.com/graph/2SAT.html

# Tree Queries

## Online Queries

that compute the result for the queries on the fly as the query comes

### Finding Ancestors

to find the kth ancestor of a given node.

to fulfill this query we compute for all n(power of two). where k<=n

and then  any k can be represented as a form of these power of 2

```

```