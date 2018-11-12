package main

import (
	"fmt"
)

type Editor interface {
	// Insert text starting from given position.
	Insert(position uint, text string) Editor

	// Delete length items from offset.
	Delete(offset, length uint) Editor

	// String returns complete representation of what a file looks
	// like after all manipulations.
	String() string
}

type Piece struct {
	origin bool // Does the piece belongs to the origin slice
	offset int  // The start of the piece
	length int  // The length of the piece
}

type Table struct {
	rows   []Piece
	origin string // the original content of the file
	add    string // the new words in the file
}

type SimpleEditor struct {
	table Table
}

func newPiece(isOrigin bool, offset, length int) Piece {
	return Piece{origin: isOrigin, offset: offset, length: length}
}

func (s *SimpleEditor) newRow(piece Piece) {

	s.table.rows = append(s.table.rows, piece)
}

func NewEditor(origin string) Editor {

	var piece Piece = newPiece(true, 0, len(origin))
	var editor SimpleEditor

	editor.newRow(piece)
	editor.table.origin = origin

	return &editor
}

// Find the index of the row in the table where the position belongs
func findPosRowIndex(position int, rows []Piece) int {

	if len(rows) <= 0 {
		return -100
	}
	var currEnd int
	resultIndex := -1
	for index, piece := range rows {

		currEnd += piece.length
		fmt.Println("The currEnd + pieceLenght are: ", currEnd)
		if position < currEnd {
			resultIndex = index
			break
		}
	}
	return resultIndex
}

// ------------------------- Insertion -------------------------

func (s *SimpleEditor) addElementAtIndex(index int, elem Piece) {
	rightPart := make([]Piece, len(s.table.rows[:index]))
	leftPart := make([]Piece, len(s.table.rows[index:]))
	copy(rightPart, s.table.rows[:index])
	copy(leftPart, s.table.rows[index:])
	if len(rightPart) == 0 {
		//fmt.Println("Add here the elem: ", elem)
		tempSlice := []Piece{elem}
		s.table.rows = append(tempSlice, s.table.rows...)
	} else {
		rightPart = append(rightPart, elem)
		s.table.rows = append(rightPart, leftPart...)
	}
	//fmt.Println("Print the rows in addElement in index function: ", s.table.rows)
}

func (s *SimpleEditor) addAtEnd(text string) {
	elem := newPiece(false, len(s.table.add), len(text))
	s.table.add += text
	s.table.rows = append(s.table.rows, elem)
}

func (s *SimpleEditor) addAtBegin(text string) {
	newRow := newPiece(false, len(s.table.add), len(text))
	tempSlice := []Piece{newRow}
	s.table.rows = append(tempSlice, s.table.rows...)
}

func (s *SimpleEditor) insertInMiddle(position, rowIndex int, text string) {

	var isOriginCurrent bool = s.table.rows[rowIndex].origin
	var currRowLength = s.table.rows[rowIndex].length
	s.table.rows[rowIndex].length = position

	//fmt.Println("Before adding: ", s.table.rows)
	newAddRow := newPiece(false, len(s.table.add), len(text))
	s.table.add += text
	s.addElementAtIndex(rowIndex+1, newAddRow)

	//fmt.Println("After adding the first element: ", s.table.rows)
	var newOffset = s.table.rows[rowIndex].length
	newSplitedRow := newPiece(isOriginCurrent, newOffset, currRowLength-newOffset)
	s.addElementAtIndex(rowIndex+2, newSplitedRow)

	//fmt.Println("The slice in the end of Insert is: ", s.table.rows)
}

// Insert text starting from given position.
func (s *SimpleEditor) Insert(position uint, text string) Editor {

	truePos := int(position)
	rowIndex := findPosRowIndex(truePos, s.table.rows)
	if rowIndex == -1 {
		rowIndex = len(s.table.rows) + 1
	}
	//fmt.Println("The position is: ", rowIndex)
	if rowIndex >= len(s.table.rows) {

		s.addAtEnd(text)
		s.table.add += text
		fmt.Println("The slice in the end of Insert is: ", s.table.rows)
	} else if truePos == 0 {

		s.addAtBegin(text)
		s.table.add += text
		fmt.Println("The slice in the end of Insert is: ", s.table.rows)
	} else {
		s.insertInMiddle(truePos, rowIndex, text)
	}
	s.deleteRowsWithZeroLength()
	return s
}

// ------------------------- Deletion -------------------------

func (s *SimpleEditor) deleteAtIndex(index int) {
	if index+1 == len(s.table.rows) {
		s.table.rows = s.table.rows[:index-1]
	} else {
		s.table.rows = append(s.table.rows[:index], s.table.rows[index+1:]...)
	}
}

func (s *SimpleEditor) deleteRowsWithZeroLength() {

	for index, piece := range s.table.rows {
		if piece.length == 0 {
			s.deleteAtIndex(index)
		}
	}
}

func (s *SimpleEditor) deleteRowsInBetween(delLength, begin, end int) int {

	rightPart := make([]Piece, len(s.table.rows[:begin]))
	leftPart := make([]Piece, len(s.table.rows[end:]))
	copy(rightPart, s.table.rows[:begin])
	copy(leftPart, s.table.rows[end:])
	s.table.rows = append(rightPart, leftPart...)

	return delLength
}

func (s *SimpleEditor) deleteInSameRow(trueOffset, trueLength, startIndex int) {
	pastLength := s.table.rows[startIndex].length
	s.table.rows[startIndex].length = trueOffset
	isOrigin := s.table.rows[startIndex].origin
	newOffset := trueOffset + trueLength
	newLength := pastLength - (trueLength + s.table.rows[startIndex].length)
	s.addElementAtIndex(startIndex+1, newPiece(isOrigin, newOffset, newLength))
}

func (s *SimpleEditor) Delete(offset, length uint) Editor {
	trueOffset := int(offset)
	trueLength := int(length)
	starDelIndex := findPosRowIndex(trueOffset, s.table.rows)
	endDelIndex := findPosRowIndex(trueLength+trueOffset, s.table.rows)

	if starDelIndex == -1 {
		return s
	} else if endDelIndex == -1 {
		s.deleteRowsInBetween(trueLength, starDelIndex, len(s.table.rows))

	} else if starDelIndex == endDelIndex {
		s.deleteInSameRow(trueOffset, trueLength, starDelIndex)
	} else {
		trueLength = s.deleteRowsInBetween(trueLength, starDelIndex, endDelIndex)
		pastEndPos := s.table.rows[starDelIndex].offset + s.table.rows[starDelIndex].length
		fmt.Println(pastEndPos)
		// Here all redundant rows are removed
		s.table.rows[starDelIndex].length -= trueLength
		s.table.rows[starDelIndex].offset = pastEndPos - s.table.rows[starDelIndex].length
	}
	fmt.Println(s.table.rows)
	s.deleteRowsWithZeroLength()
	return s
}

// ------------------------- String -------------------------

// String returns complete representation of what a file looks
// like after all manipulations.
func (s *SimpleEditor) String() string {

	var result string
	fmt.Println("The editor rows are: ", s.table.rows)
	fmt.Println("The origin array: ", s.table.origin)
	fmt.Println("The add array: ", s.table.add)
	for _, piece := range s.table.rows {
		offsetIndex := piece.offset
		length := piece.length
		fmt.Println("Piece is: ", piece)
		if piece.origin {
			result += s.table.origin[offsetIndex : length+offsetIndex]
		} else {
			result += s.table.add[offsetIndex : length+offsetIndex]
		}
	}
	return result
}

func main() {

	// var editor SimpleEditor = NewEditor("A ")
	// fmt.Println(len(editor.table.rows))
	// fmt.Println(editor.table.origin)

	// fmt.Println(editor.String())

	//editor := NewEditor("A large span of text")

	//var editor SimpleEditor
	//fmt.Println(editor.table.rows)
	//editor.newRow(newPiece(true, 0, 8))
	// editor.newRow(newPiece(false, 0, 4))
	//editor.table.add += "English "

	//editor = editor.Insert(16, "English ")

	//fmt.Println(editor.String())

	//fmt.Println(editor.table.rows)

	//editor.Insert(7, "B")

	// fmt.Println(editor.String())
	// fmt.Println(editor.table.rows)

	//editor.Insert(0, "Nooo way ")

	// fmt.Println(editor.String())
	// fmt.Println(editor.table.rows)

	// editor.Delete(2, 6)

	// fmt.Println(editor.String())

	// editor.Delete(0, 9)

	// fmt.Println(editor.String())

	// editor.Delete(10, 150)

	// fmt.Println(editor.String())

	// editor := NewEditor("Some text")
	// fmt.Println(editor.String())
	// editor = editor.Insert(4, " random")
	// fmt.Println(editor.String())
	// editor = editor.Insert(15, " here")

	editor := NewEditor("Some random text")
	fmt.Println(editor.String())
	editor = editor.Delete(11, 4)

	fmt.Println(editor.String())
	fmt.Println("--------------------------------------------------------------")

	// editor.newRow(newPiece(true, 0, 16))
	// editor.newRow(newPiece(false, 0, 8))
	// editor.newRow(newPiece(true, 16, 4))

	// fmt.Println(editor.table.rows)

	// editor.newRow(newPiece(false, 8, 1))

	// editor.addElementAtIndex(0, newPiece(false, 100, 100))

	// fmt.Println(editor.table.rows)

	//fmt.Println(findPosRowIndex(21, editor.table.rows))

	// var editor SimpleEditor

	// editor.newRow(newPiece(true, 16, 4))

	// fmt.Println(findPosRowIndex(16, editor.table.rows))

	//var editor = NewEditor("A large span of text")

	//fmt.Println(reflect.TypeOf(editor.Insert(16, "English ")))

	//fmt.Println(reflect.TypeOf(editor))

	//fmt.Println(editor.String())

	//editor.Insert(30, " cool")

	//fmt.Println(editor)

	// editor.Insert(300, " daaa be")

	// fmt.Println(editor)

	//editor.Insert(4, "1")
	//fmt.Println(editor)

	// var try := "A large span of text"

	// fmt.Println(try[16:20])
	//editor.Insert(16, "English ")
	//editor = SimpleEditor(temp)

	//fmt.Println(editor.String())

	fmt.Println("Finished")
	//fmt.Println(temp.String())

	// slice := []int{0, 1, 100, 4}

	// x := append(slice[:2], slice[2+1:]...)
	// fmt.Println(x)

	// slice = addAtIndex(slice, 0, 13)

	// slice = addAtIndex(slice, 1, 100)

	// slice = addAtIndex(slice, 4, 999)
	//right := []int{13, 100}
	//left := []int{0, 1}

	//fmt.Println(append(right, left...))

}
