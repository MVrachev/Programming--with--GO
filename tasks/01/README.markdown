# Poor Man's Currying

In this task you will create a few functions which generate other functions. This is very similuar to [currying](https://en.wikipedia.org/wiki/Currying).

## Repeater

Function that accepts as arguments string `s` and separator `sep` and returns a function which concatenates the string given times using the string  `sep` as separator and returns the result string.

```
func Repeater(s, sep string) func (int) string
```

Sample use:

```
Repeater("foo", ":")(3) // foo:foo:foo
```

## Generator

Function which creates a "generator" function for `int` numbers.

```
func Generator(gen func (int) int, initial int) func() int
```
You pass a `gen` function to Generator and start value in the sequence `initial`. `gen` it takes as first argument previous calculated value and returns the next.

Sample use:

```
counter := Generator(
    func (v int) int { return v + 1 },
    0,
)
power := Generator(
    func (v int) int { return v * v },
    2,
)

counter() // 0
counter() // 1
power() // 2
power() // 4
counter() // 2
power() // 16
counter() // 3
power() // 256
```

## MapReducer

Function, which creates map reducer function for `int` arguments passesd [map](https://en.wikipedia.org/wiki/Map_(higher-order_function)) function, [reduce](https://en.wikipedia.org/wiki/Fold_(higher-order_function)) function and the inital value 
`initial` for the reduce function.

```
func MapReducer(mapper func (int) int, reducer func (int, int) int, initial int) func (...int) int
```

You can call MapReducer this way:

```
powerSum := MapReducer(
    func (v int) int { return v * v },
    func (a, v int) int { return a + v },
    0,
)

powerSum(1, 2, 3, 4) // 30
```

The arguments of the `reducer` function should be passed from left to the right or with other words create a [left-fold](https://en.wikipedia.org/wiki/Fold_(higher-order_function)#On_lists) implementation of the reduce function.

## Reminder

Don't forget to [format your code using gofmt](https://blog.golang.org/go-fmt-your-code).

