# B+ Tree Implementation in Go

This repository contains an implementation of a **B+ Tree** data structure written in Go. A B+ Tree is a self-balancing tree commonly used in databases and file systems due to its efficient support for range queries, sequential access, and balanced structure. This implementation supports insertion, deletion, modification, and search operations, with additional utilities for debugging and visualization.

## Features

- **Dynamic Order**: Configurable maximum number of keys per node (`MaxKeys`), with automatic calculation of minimum keys.
- **Leaf-Linked Structure**: Leaf nodes are linked via a `next` pointer, enabling efficient sequential traversal.
- **Insertion**: Handles node splitting for both leaf and internal nodes when exceeding the maximum key limit.
- **Deletion**: Supports rebalancing through borrowing from siblings or merging nodes to maintain the minimum key requirement.
- **Search**: Efficiently locates a key and returns its associated value, or `-1` if not found.
- **Modification**: Updates the value associated with an existing key.
- **Debugging Tools**: Includes functions to print the tree structure and leaf node values for easy visualization.

## Prerequisites

- **Go**: Version 1.13 or higher (due to the use of modern slice manipulation and error handling).

## Installation

1. Clone this repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Ensure you have Go installed. Verify with:
   ```bash
   go version
   ```

3. Run the example:
   ```bash
   go run main.go
   ```

## Usage

The B+ Tree is implemented as a struct `BPlusTree` with methods for core operations. Below is a basic example of how to use it:

```go
package main

import "fmt"

func main() {
    tree := NewBPlusTree()

    // Insert key-value pairs
    tree.Insert(1, 10)
    tree.Insert(2, 20)
    tree.Insert(3, 30)

    // Search for a key
    value := tree.Search(2)
    fmt.Println("Value for key 2:", value) // Output: 20

    // Modify a value
    tree.Modify(2, 25)
    fmt.Println("New value for key 2:", tree.Search(2)) // Output: 25

    // Remove a key
    tree.Remove(1)
    fmt.Println("Value for key 1 after removal:", tree.Search(1)) // Output: -1

    // Print the tree structure
    tree.PrintTree()
}
```

### Configuration

- **MaxKeys**: The maximum number of keys per node is defined as a constant (`MaxKeys = 3` by default). Adjust this value in the code to change the tree's order.
- **Minimum Keys**: Automatically calculated based on `MaxKeys` (e.g., `MaxKeys / 2` for even numbers, `(MaxKeys + 1) / 2` for odd).

## Code Structure

- **`Node` Struct**: Represents a node in the B+ Tree.
  - `isLeaf`: Boolean indicating if the node is a leaf.
  - `keys`: Slice of integers storing keys.
  - `values`: Slice of integers storing values (leaf nodes only).
  - `children`: Slice of pointers to child nodes (internal nodes only).
  - `next`: Pointer to the next leaf node (leaf nodes only).
  - `parent`: Pointer to the parent node.

- **`BPlusTree` Struct**: Represents the B+ Tree itself, with a single field `root`.

- **Key Methods**:
  - `Insert(key, value int)`: Inserts a key-value pair.
  - `Remove(key int) error`: Deletes a key and rebalances the tree if needed.
  - `Search(key int) int`: Searches for a key and returns its value.
  - `Modify(key, newValue int) error`: Updates the value of an existing key.
  - `PrintTree()`: Prints the tree structure level by level.
  - `PrintLeafValues()`: Prints all values stored in leaf nodes sequentially.

- **Helper Functions**:
  - `splitLeaf` and `splitInternal`: Handle node splitting.
  - `rebalance`: Ensures nodes meet the minimum key requirement after deletion.
  - `findLeaf`: Locates the appropriate leaf node for a given key.
  - `updateInternalKeys` and `updateParent`: Maintain consistency of internal node keys.

## Example Output

Running the `main.go` file produces output demonstrating insertion, deletion, and search operations:

```
插入数据：
插入 (4, 40)
插入 (7, 70)
插入 (1, 10)
插入 (9, 90)
插入 (2, 20)
插入 (5, 50)
插入 (8, 80)
插入 (3, 30)
插入 (6, 60)
初始B+树结构：
[Internal, -1: 3 6 9 ]
[Leaf, 3: 1 2 3 ]  [Leaf, 6: 4 5 6 ]  [Leaf, 9: 7 8 9 ]
所有叶节点对应的值：10 20 30 40 50 60 70 80 90

删除 keys: 9, 1, 5, 2
删除 9 后的结构：
[Internal, -1: 3 6 8 ]
[Leaf, 3: 1 2 3 ]  [Leaf, 6: 4 5 6 ]  [Leaf, 8: 7 8 ]
所有叶节点对应的值：10 20 30 40 50 60 70 80

删除 1 后的结构：
[Internal, -1: 3 6 8 ]
[Leaf, 3: 2 3 ]  [Leaf, 6: 4 5 6 ]  [Leaf, 8: 7 8 ]
所有叶节点对应的值：20 30 40 50 60 70 80

删除 5 后的结构：
[Internal, -1: 3 6 8 ]
[Leaf, 3: 2 3 ]  [Leaf, 6: 4 6 ]  [Leaf, 8: 7 8 ]
所有叶节点对应的值：20 30 40 60 70 80

删除 2 后的结构：
[Internal, -1: 6 8 ]
[Leaf, 6: 3 4 6 ]  [Leaf, 8: 7 8 ]
所有叶节点对应的值：30 40 60 70 80

查找 key 7 对应的 value: 70
```

## Notes

- **Memory Management**: The implementation relies on Go's garbage collector to handle memory deallocation, avoiding manual freeing of nodes.
- **Error Handling**: Deletion and modification operations return errors if the key is not found.
- **Thread Safety**: This implementation is not thread-safe. For concurrent use, add synchronization mechanisms (e.g., mutexes).

## Contributing

Contributions are welcome! Please submit a pull request or open an issue for bug reports, feature requests, or improvements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
