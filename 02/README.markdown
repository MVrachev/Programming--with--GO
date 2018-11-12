# Piece table

In this task you will have to implement the main functionality of one text editor.
He has to work on the on the potential content of a file and give the ability to inser, delete text and undo and redo operations.

One famous structure for implementation of this funcionality is [piece
table](https://en.wikipedia.org/wiki/Piece_table) (in some places you can find it as a piece chain).

For the implementation you need two buffers:

- First one (file buffer) contains the original content of the file. The file buffer is immutable (read-only). In this task we will call it "origin".

- The second buffer contains only the additional content of the file added by insertion. The content of this buffer can only grow, once data is added there it's never deleted (append-only). In this task we will call it "add".

The table will contains "pieces" which will relate to a part of the file based on that in which buffer they are located.
We can present one piece with the following structure:


	struct {
		origin bool
		offset int
		length int
	}

- origin says does the content of the piece belongs to the origin buffer.

- offset gives us the from which **byte** in the в aforenamed buffer starts the content of that piece

- length gives us the length of the piece.

In the begging when opening a file the whole fyle is saved in one piece, which is the only record in the table.

### Example:

We open a file with the content `A large span of text`. The origin buffer looks like this:

![origin](./images/origin.png)

The add buffer is empty and the table looks like this:

![table0](./images/table0.png)

The user decides to add the word "English" before the last word to get a file with the content: `A large span of English text`. We add `English ` in the add buffer.

![add](./images/add.png)

Note that the size of the add buffer is larger than the necessary. That's okay but it's not mandatory. 
The table now looks like this: 

![table1](./images/table1.png)

Our user decides to delete the word `large`. This doesn't lead to any changes in our buffers and doesn't delete anything from them because we can't delete from the buffers.

The table changes into this: 

![table2](./images/table2.png)

## The task

Your task is to implement a type which implements the following interface:


	type Editor interface {
		// Insert text starting from given position.
		Insert(position uint, text string) Editor

		// Delete length items from offset.
		Delete(offset, length uint) Editor

		// String returns complete representation of what a file looks
		// like after all manipulations.
		String() string
	}

Also you are expected to create a function, which gives value to your own type from the given origin buffer.


	func NewEditor(string) Editor

### An example

The example before we can reproduce with your implementation like this: 


	var f = NewEditor("A large span of text")
	f.String() // "A large span of text"

	f = f.Insert(16, "English ")
	f.String() // "A large span of English text"

	f = f.Delete(2, 6)
	f.String() // "A span of English text"

## Notes / Recommendations:

- Note that the functions doesn't return an error. That means that we expect from you to handle with too big values of position, offset and length:
    - Too big value of position in Insert should act like a normal Insert into the end of the file: 
      `NewEditor("foo").Insert(453, ".")` should work as:
      `NewEditor("foo").Insert(3, ".")`.
    - Too big value of the offest in Delete shouldn't change the file. 
      This operation doesn't do anything:
      `NewEditor("foo").Delete(300, 1)`
    - Too big value of length in Delete should go to the end of the file
      `NewEditor("foo").Delete(1, 300)` should act like this:
      `NewEditor("foo").Delete(1, 2)`.
- Note that every method returns an editor. This still doesn't mean that is mandatory to create a new editor which you will return. The descision is yours to change the existing editor or to create a new one in every operation. 
- add can grow infinitely.
- Don't focus to optimize the allocation too much. 
- When opening a new few its okay the add buffer to be nul
- Your type **should not** be caled `Editor`, because you already have defined one in the sample tests.
- The Insert and Delete operation expects the position, offset and length in **bytes** not in symbols. Note this when working with unicode symbols.
