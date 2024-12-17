#notes #books #compititive #programming #ibasicTechinques


# Sorting Algorithms
#sorting #algos

## O(n^2) 

Bubble Sort 
```
n rounds of sorting and swapping the two numbers which are not in the right location
```

## O(nlogn)

Merge Sort 

```
1. calculate the middle element
2. recursively sort the [a..k]. [k+1.....b]
3. now merge the sorted element to reveal a sorted [a, b]
```

The lowest bound for sorting algos is O(nlogn), we cant go any lower that . if we represent the choices taken at each comparison as tree, the minimum number of decision to reach the final correct outcome , is one at each level, and hence the height of a binary tree, nlogn


# Binary Search
#binarySearch 

Binary searches in a sorted array a given index in O(logn)

Two methods to implement the binary search algo

Method 1
```cpp
int a=0,b=n-1;
while(a<=b){
	int mid = (a+b)/2;
	if(array[mid]==target){
		//do something
	}else if(array[mid]>target) b=mid-1;
	else a=mid+1;
}
```

Method 2
This method on jumping through the false condition by reducing the jump by half on every traversal

```cpp
int start=0;
for(int b=n/2;b>=1;b/=2){
	while(start+b<n && array[start+b]==target) start+=b;
}

if(array[start]==target)
//do something
```

Examples of the use of method2

Example 1
we need to find a change in the function, the k at which the values switched in the array

```cpp
int x=-1;
for(int b=z;b>=1;b\=2){
	while(!func(x+b)) x+=b;
}

int k = x+1;
```

this finds the largest value of x for which the func returns false;
thus the x+1 is the first value when the func returns `true`;

Example 2
our array or function is first increasing and then decreasing, we need to find the first point it starts decreasing or the inflection point

```cpp
int x=-1;
for(int b=z;b>=1;b\=2){
	 while(f(x+b) < f(x+b+1)) x+=b;
}

int k = x+1;
```

z is the first jump that we know the condition in the while loop stands true

# Data Structures
#dataStructures

```
1. vector<>
2. set<>
3. unordered_set<>
4. multiset<>
5. map<>
6. unordered_map<>
7. bitset<10>
8. deque<>
9. queue<>
10. priority_queue<>
11. stack<>
```

# Complete Search
#completeSearch

```
the idea is to get all the possible solutions to a problem
using brute force or otherwise
```

## Generating Subset

Method 1
gets the recursion call with and without the current n added in the subset
```cpp
void search(int k){
if(k==n){
	//process
}else{
	search(k+1);
	subset.push_back(k);
	search(k+1);
	subset.pop_back();
}
}
```

Method 2
with the bit manipulation method

its an intuitive method where for lets say a given number n for which we have to find the subsets, we take 2^n-1 .
then we try to check from 000 to 2^n-1 we check for the set bits,
 if the bits are set then the number is included in the subset, rightmost bit is 0, then 1, then 2 and so on

```cpp

for(int b=0;b<(1<<n);b++){
	vector<int> subset;
	for(int i=0;i<n;i++){
		if(b&(1<<i)) subset.push_back(0);
	}
}
```

## Backtracking
#backtracking

a form of complete search, we start with an empty solution and recursively try to get the best possible solution in each scenario

lets say the question is to place the n queens on a board so they dont attack each other
We should start with placing the queen in the first row
and the second row we can fin the appropriate place we are not attacked with the queen in the first row.

```cpp
void search(int y){
	if(y==n){
		count++;
		return;
	}
	for(int x=0;x<n;x++){
		// dont add a queen to the position whose column or the diaganal    already contains a queen
		if(column[x] || diag1[x+y] |\ diag[x-y+n-1]) continue;
		column[x] = diag1[x+y] = diag2[x-y+n-1] = 1; 
		search(y+1); 
		//release the column and the diagnal so we can get other open positions as well
		column[x] = diag1[x+y] = diag2[x-y+n-1] = 0;
	}
}
```


## Pruning the search

clearly the above method is really time consuming and would take huge times for bigger matrices.
hence the only logical solutions is to prune the search areas

these optimizations are based on the problems statement istle
sometimes its better to divide the problem subset into way which significantly reduces the overall load  as well.

# Greedy Algos
#greedy

```
always making a choice that is the  best at that moment, so it doesnt take back the choice but constructs the final solution.

so the algo has to ensure that the locally optimal solution is also globally optimal

so many a time the greedy approach is not always optimal and might not give the right solution
```

What works for the greedy approach ?

1. Scheduling / Tasks / Deadlines
```
	when you have the schedule the maximum possible tasks / events within  a given time, in these case finding the best optimal solution at each point, works the best.
	
	So lets you can select the event which ends the earliest, which will give the time to select the other events
	Or given a set of tasks and deadlines, select the order to perform those tasks in. its beter to complete the lower duration tasks quickly at each step
```

2. Data Compression / Huffman Coding
```
Huffman is also a greedy algorithm which first takes the occurences of each character in the string.
and at each step the optimal; solution is to get the least occuring two characters and create a binary tree of them with their sum as the root node

this way you can construct a codeword for each occurence in the string 
and then replace the string with the same occurencce to compress it

```

# Dynamic Programming
#dp

```
DP essentially combines the complete search and the greedy algo.

DP can be used for
1. find an optimal solution, min steps to reach
2. couting the number of solutions
3. longest increasing subsequence
4. Paths in a grid
5. Knapsack problem
6. Edit Distance
7. Counting Tilings


with DP essentially the answer is through caluclating and dividing the task into subproblems

to do this its usually a combination of the recursion, calculating the subproblems + memoization , storing the solutions of already caclculated subproblems.

because recursion calls and overall no of paths can be exponential.
storing the the already calculating solution reduces the task to only go through every possible path once


```


# Range Queries
#rangeQueries

## Range Sum Query

these can be done easily using the sum prefix table, which can be computed in O(n)
time and then the the range(a,b) = prefix(b)-prefix(a) to get the answer in O(1)


## Minimum / Maximum / GCD Range Query

what if we need to find the minimum no in a given range, over a list of queries

we need to precompute for every range which is O(n^2) might not work in every scenario.

How to precompute properly

1. Sparse Table

sparse table help compute time to logarithmic, to cmpute for only power of 2
and hence the precompute time is O(nlogn) and the query time is O(1)
//for precompute
![[../Images/Pasted image 20241209175339.png]]


//for query minimum
![[../Images/Pasted image 20241209173716.png]]

so for the precomputation we compute it for 2^ (0->K), where k is the log2(length of the range)

why the log2 ? because it just  means how many times to divide from 2 to get till 1, so gives us the maximum power less than the length

once its computed for every power, we store it a matrix for every i as row, the j (power) as the column

we can use the query to get the required info.
again  let say given query is [a,b] 

```
length of the range := b-a+1
j = log2(b-a+1)
min(a,b) =min(precompute[a][j], precompute[b-(1<<j)+1][j])
```

the code, can refer to here
https://www.youtube.com/watch?v=0jWeUdxrGm4

```cpp
//actual array
int arr[n]

//computes the log of the number for quick access
int logcompute[]

//log 2 just means keep dividing until 1
for(int i=2;i<=n;i++)
	logcompute[i]=logcompute[i/2]+1;


//precompute
int precompute[n][j];

//compute the lenght 1, j==0 which will be same the value on the index
for(int i=0;i<n;i++){
	precompute[i][0]=arr[i];
}

//compute for the rest j's
for(int j=1;j<logcompute[n];j++){
	for(int i=0;i+(1<<j)-1<n;i++){
		precompute[i][j] =min(precompute[i][j-1], 
							  precompute[i+(1<<j)][j-1]);
	}
}

//for a given query [a,b]
int len = b-a+1
int j = logcompute[len];
min(a,b)=min(precompute[a][j], precompute[b-(1<<j)+1][j]);
```


this can applied for the maximum and the gcd as well

## Binary Indexed Tree or the Fenwick Tree

this is modification of the range sum query, the addition being that it allows the updating a value, which in the case of the prefix would have been recomputing everything
both the processing and updating the value can be done in O(logn) times.

This is based on the previous range query to find the maximum/minimum

where we depend on power 2 to comupte the ans for value possible ranges
in this case we compute the sum the for the power of 2

https://www.youtube.com/watch?v=uSFzHCZ4E-8

![[../Images/Pasted image 20241209185559.png]]

we compute the tree(its an array, called tree). we copy the same vale if 0th is set to 1,
otherwise for every other set bit we sum 26 elements
for examples 00010, has 2 elements
00100, has 4 elements
01000, has 8 elements


How to find the sum at index i ?
we keep on flipping the last bit to get the information stored in the already computed previous Tree Array
![[../Images/Pasted image 20241209190507.png]]
```cpp
//to find the sum at a particular index

int sum(int i){
int sum=0;
while(i>0){
	sum+=T[i]
	i-=i&-i; // flip the last set bit
}
return sum
}

```


what does the i-=i&-i do

```
7= 00111
-7= 11000 + 1 = 11001 (2's complement)

00111 && 11001 = 00001
00111 - 0001 = 00110
```


Now we want to update an index, we want to update all the computed sum the index was part , because we logN computation,  logn updates are required

```cpp

int add(int i , int k){
while(i<tree.length()){
	tree[i]+=k;
	k+=k&-k; //add last set bit
}
}
```


The actual computaion of the tree array can also be done similarly, adding the current index to its parent range

```cpp

int[] make(int[] arr){
	int[] tree = Arrays.copyof(arr);
	for(int i=1;i<n;i++){
		int nextParent = i+(i&-i);
		if (nextParent < n){
			tree[parent]+=tree[i];
		}
	}
}
```

## Segment Tree

segment trees is sort of a beefy version of the fenwick tree, takes 2n storage space.
But calculates all the range queries in logn, including sum / maximum / minimum.

It constructs a tree out of of every two indexes, propogating it upwards
![[../Images/Pasted image 20241210152630.png]]


How to construct the array space - O(2n) , time :- O(n)

high = n-1, low = 0 , size of segtree is 2*n, pos=0
```cpp
void contructType(int[] input, int segTree[], int low, int high, int pos){
if (low==high){
	segTree[pos] = input[low];
	return ;
}

int mid = (high+low)/2;
contructTree(input, segTree, low, mid, 2*pos+1) //left child
contructTree(input, segTree, mid+1, high, 2*pos+2)//right child
segTree[pos]= min(segTree[2*n+1], segTree[2*n+2]);
}
```

How to query [a,b]

```cpp
int query(int ind, int low, int high, int a , int b){
	if(high<=b && low>=a){ // lies completely in the range
		return seg[ind];
	}
	if(high<a || low>b){ //lies completely out of the current range og array
		return INT_MAX;
	}
	int mid = (low+high)/2;
	int left = query(2*ind+1,low, mid ,a,b);
	int right = query(2*ind+2,mid+1, high ,a,b);
	return min(left, right);
}
```

So the root level node in the tree, denotes the whole range of the array, and then the child nodes denotes the division in the middle.

so if the top level node is [0,9]
the left node would be[0,4]
the right node would be [5, 9]

for the given range [a,b] we try to find if the given node's range is within a,b
otherwise we want to go both left and right nodes

Segment trees also support the index updation in the O(logn) time

because we trace the path or the node where the index update would affect,
basically if its in the node range or not


```cpp
int updateNode(int ind, int low , int high, int index, int val){
if(low==high){ // the root node
	seg[low]+=val;
	return;
}

int mid =(high+low)/2;
if(index<=mid && index>=low){ //in the left subtree
	updateNode(2*ind+1, low, mid, index, val);
}else{ // in the right subtree
	updateNode(2*ind+2, mid+1, high, index, val);
}

seg[ind]=seg[2*ind+1]+seg[2*ind+2];

}
```


lazy propagation

for a better understanding, but basically, given a range [a,b] you wanna make an update to the all of the indexed within that array.

Now with lazy propogation, you only update till the node which completely lies between the given [a,b]
And store the update for its children, down the line when another query comes you can take this pending update + actual query update and update the node.

this way you only update what is required


```cpp
int[] lazy ={0}, same size as segment tree
int updateRange(int ind, int low, int high, int a, int b, int val){
	if(lazy[ind]!=0) {//pending updates from another query
		seg[ind] += (high-low+1) * lazy[ind]// the range [0.2] has three node and take update for eah
		if(low!=high){//if its not the leaf node, update the lazy for  children
			lazy[2*ind+1]+=lazy[ind];
			lazy[2*ind+2]+=lazy[ind];
		}
		lazy[ind]=0;
	}

	if(b<low || a>high || low>high) return; //out of range

	if(low>=a && high<=b){//lies completely in the range
		seg[ind]+=(high-low+1)*val;
		if(low!=high){
			lazy[2*ind+1]+=val; //lazy updatte the children node for other queries
			lazy[2*ind+2]+=val;
		}
	return;
	}

//range  lies partially in both left and right subtrees
int mid=(high+low)/2;
updateRange(ind, low, mid, a, b, val);
updateRange(ind, mid+1, high, a, b, val);

seg[ind]=seg[2*ind+1]+seg[2*ind+2];
}
```

similar way to get the sum of the nodes after an update, we can check the lazy array for any pending updates, and hence use those updates+ propagates them further for the next queries

# Bit Manipulation
#bitManipulation

unsigned int0 2^n -1
signed int -2^n-1 to 2^n-1 -1

## Bit Operations

AND operation
```
x & 1 == 0 // even
X & 1 == 1 // odd
x & (2^k -1) ==0 // divisible by 2^k
```

OR operation a | b
XOR operation a ^ b

NOT operation
```
~x = -x-1
```

BIT shifts

```
left shift  x<<k , appends k 0 bits, or multiply by 2^k
right shift x>>k  removes k 0 bits, o r divde by 2^k

eg 14<<2=56
14=1110
56=111000
```


## Applications

1. `x & (1<<k)` , checks if the kth bit of x is set or not
2. `x | (1<<k)` , sets the kth bit in x to 1
3. `x ^ (1<<k)`, flips the kth bit in x
4. `x & ~(1<<k)`, set the kth bit to 0
5. `x & (x-1)`  set the last 1  bit of x to 0
	1. x = 7, 111 ; x-1 = 6, 110; result = 110
	2. x=6, 110; x-1 = 5, 101; result = 100
6. `x & -x`, sets all 1 bits to 0, except last 1 bit
	1. x=7, 111 ; -x=000+1=001; result = 001
	2. x=6, 110; x=001+1=010; result =010
7. `x | (x-1)` inverts all bits after last 1 bit
	1. x=7, 111; x-1=6, 110; result =111
	2. x=6, 110; x-1 = 5, 101; result 111

Built In Functions
`__builtin_clz()` number of 0's at the beginning of the number
`__builtin_ctz()` number of 0's at the end of the number
`__builtin_popcount()` number of 1's in the number
`__builtin_parity()` parity of number of 1's
