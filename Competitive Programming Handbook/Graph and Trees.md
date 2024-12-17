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

### Finding Ancestors

to find the kth ancestor of a given node.

Brute Force -> find the node and then recursively get the parent
This works well for balanced Tree

Precomputation

to fulfill this query we compute for all n(power of two). where k<=n

and then  any k can be represented as a form of these power of 2
in a way precomputes the recursive calls for a large number of queries to be  processed

Works best for unbalanced trees


```cpp
vector<vector<int> up;
int LOG;
void precomputation(int n, vector<int> parent){
     LOG=ceil(log2(n));
     precompute=vector<vector<int>>(n, vector<int>(LOG+1,-1));

    for(int v=0;v<n;v++){
		precompute[v][0]=parent[v];
	}

	 for(int j=1;j<=LOG;j++){
		for(int v=0;v<n;v++){
			int nextAncestor =precompute[v][j-1];
			if(nextAncestor!=-1){
				precompute[v][j]=precompute[nextAncestor][j-1];
			}
		}

}

```


now to calculate the kth ancestor, we can find the bits that are set in the k to get that ancestor

```cpp

    int getKthAncestor(int node, int k) {
        for(int j=0;j<=LOG;j++){
            if(k & (1<<j)){
                node = precompute[node][j];
                if(node==-1)
                    return -1;
            }
        }
        return node;
    }
};

```


### Subtree Queries / Path Queries

we can denote the subtree dfs approach in the form of an array
if we given the size of subtree of the node, we can use the range query to

update the value of a node
sum of the subtree of the given node

hence we construct a array with 3 values in each node
node id, subtree size,  node value


![[Pasted image 20241216172211.png]]

here we calculate sum of subtree of node 4

![[Pasted image 20241216172233.png]]


Similarly we can get the path sum using the similar approach, where each path sum we keep on adding each node while we dfs
![[Pasted image 20241216172520.png]]

### Lowest Common Ancestor

Method 1 - Simple DFS
#### Method 2 - Binary Lifting

precompute the depth and the log parent of each node as in the previosu examples

then tow find the LCA, we need move A and B until there parents are not equal and then return the parent.

Which makes sense because  when the condition breaks and we are at a point where parents are equal , which means it is the lowest common ancestor

```cpp
int lca(int nodeA, int nodeB){
	if(depth(nodeA)<depth(nodeB)){
		swap(nodeA, nodeB);
	}
	
	int k = depth(nodeA)-depth(nodeb);
	
	//bring both the node on the same level
	for(int j=0;j<=LOG;j++){
		if(k & (i<<j)){
			nodeA = precompute[nodeA][j];
		}
	}
	//if both the nodes are same, then return the node
	if(a==b)
		return a;
    //otherwise we move till we get the best option
    for(int j=0;j<=LOG;j++){
	    if(precompute[nodeA][j]!=precompute[nodeB][j]){
		    nodeA=precompute[nodeA][j];
		    nodeB=precompute[nodeB][j];
	    }
    }

	//the parent of the node when the condition breaks is the answer
	return precompute[nodeA][0];
}
```


#### Method 3 - Tree Traversal / Euler  Tour Technique

![[Pasted image 20241216183615.png]]

![[Pasted image 20241216183631.png]]

here while dfs, we add node every time, we come across it in the dfs.
(while going in depth and coming back up)

then range minimum queries can be used to find the ancestor with the least depth



