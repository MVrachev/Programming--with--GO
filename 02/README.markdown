# Piece table

In this task you will have to implement the main functionality of one text editor.
He has to work on the on the potential content of a file and give the ability to inser, delete text and undo and redo operations.

One famous structure for implementation of this funcionality is piece
table](https://en.wikipedia.org/wiki/Piece_table) (in some places you can find it as a piece chain).

For the implementation you need two buffers:

- First one (file buffer) 

За имплементацията на таблицата, са ни нужни два буфера:

- Първият (file buffer) съдържа целият зареден файл, преди всякакви редакции.
  Той никога не търпи промени в процеса на работа (read-only). В рамките на
  тази задача ще го наричаме "origin".

- Вторият буфер съдържа само добавеното съдържание във файла. Съдържанието на
  този буфер може само да расте, веднъж добавени данни там, те не търпят
  изтриване в процеса на работа (append-only). В рамките на тази задача ще
  го наричаме "add".

Таблицата съдържа "парчета", които описват част от съдържание във файла, въз
основа на това, в кой буфер се намират. Едно парче можем да представим със
следната структура:

	struct {
		origin bool
		offset int
		length int
	}

- origin указва дали съдържанието, към което сочи се намира в origin буфера.

- offset указва от кой **байт** в горепосочения буфер започва съдържанието на
  това парче.

- length указва дължината на въпросното съдържание.

При първоначалното отваряне на даден файл целият файл се описва от едно парче,
което е единственият запис в таблицата.

### Пример:

Отваряме файл със съдържание `A large span of text`. Origin буферът изглежда така:

![origin](./images/origin.png)

Add буферът е празен, а таблицата с парчета изглежда така:


![table0](./images/table0.png)

Потребителят решава да добави думата "English" преди последната дума, за да
получи файл със съдържание: `A large span of English text`. Добавяме в add буфера `English `:

![add](./images/add.png)

Обърнете внимание, че  размерът на add буфера е по-голям от необходимото. Това
е допустимо, но не задължително. Таблицата ни вече изглежда така:

![table1](./images/table1.png)

Въображаемият ни потребител решава да изтрие думата `large`. Това не води до
никакви промени в нашите буфери, тъй като нищо никога не се трие от тях.
Таблицата се променя до:

![table2](./images/table2.png)

## Задача

Задачата ви е да създадете тип, който имплементира следния интерфейс:

	type Editor interface {
		// Insert text starting from given position.
		Insert(position uint, text string) Editor

		// Delete length items from offset.
		Delete(offset, length uint) Editor

		// Undo reverts latest change.
		Undo() Editor

		// Redo re-applies latest undone change.
		Redo() Editor

		// String returns complete representation of what a file looks
		// like after all manipulations.
		String() string
	}

Също така се очаква и да добавите функция, която създава стойност на вашия тип,
от подаден origin буфер.

	func NewEditor(string) Editor

### Пример

Горният пример би трябвало да може да бъде пресъздаден с помощта на вашата
имплементация по този начин:

	var f = NewEditor("A large span of text")
	f.String() // "A large span of text"

	f = f.Insert(16, "English ")
	f.String() // "A large span of English text"

	f = f.Delete(2, 6)
	f.String() // "A span of English text"

## Забележки/Препоръки:

- *Внимание*: когато предавате решение, трябва да включите в него дефиницията
  на `Editor` интерфейса, за да успеете да предадете решението си.
- Обърнете внимание, че никой от методите не връща грешка. Това ще рече, че
  очакваме да се справяте с твърде големи стойности на position, offset и
  length:
    - Твърде голяма стойност на position при Insert трябва да се държи като
      обикновено добавяне на края на файла.
      `NewEditor("foo").Insert(453, ".")` трябва да се държи като
      `NewEditor("foo").Insert(3, ".")`.
    - Твърде голяма стойност на offset при Delete не трябва променя файла.
      Тази операция не прави нищо:
      `NewEditor("foo").Delete(300, 1)`
    - Твърде голяма стойност на length при Delete трябва да стига до края
      на файла. `NewEditor("foo").Delete(1, 300)` трябва да се държи като
      `NewEditor("foo").Delete(1, 2)`.
- Обърнете внимание на факта, че всеки един от методите връща редактор. Това
  все пак **не** означава, че е задължително да създавате нов редактор, който
  да връщате. Решението дали да го правите, или да мутирате текущия е ваше.
- add буферът може да расте до безкрайност.
- Не се фокусирайте да оптимизирате тези алокации особено много. Напълно ок е
  да решите, че по подразбиране add буферът е X символа и ако се налага просто
  удвоява размера си.
- При отваряне на файл add буферът е празен и е напълно в реда на нещата дори
  да не е алокиран (т.е. да е nil).
- Вашият тип **не трябва** да се казва `Editor`, тъй като вече имате такъв
  дефиниран интерфейс в sample тестовете. Напълно допустимо е този тип да не
  бъде exported.
- Undo преди извършена операция не трябва да прави нищо.
- Redo преди изпълнено Undo (или след всички възможни Redo-та) не трябва да
  прави нищо.
- След Undo всяка редакция се предполага да инвалидира следваща Redo операция.
- Операциите Insert и Delete очакват position, offset и length в **байтове**, а
  не в символи. Обърнете внимание на това при работа с unicode символи.