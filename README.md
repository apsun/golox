# golox

This is my attempt at following along to the excellent
[Crafting Interpreters](https://www.craftinginterpreters.com/) book.

It's written in Go because I wanted to get better at Go (and because I'm a
masochist). For the most part, it's just a direct translation of the Java
code from the book into Go, with some slight design adjustments to be more
idiomatic.

WARNING: I decided to skip the AST visitor pattern out of laziness, and
instead crammed all of the logic within the Expr/Stmt types themselves.
In hindsight, this was an awful idea and has led to countless hours wasted
untangling spaghetti code; if you are using this code for inspiration, do
not copy this particular design choice.
