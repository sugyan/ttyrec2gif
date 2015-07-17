ttyrec2gif
==========

ttyrec to animated GIF.

![](https://cloud.githubusercontent.com/assets/80381/3735009/5fa96246-171f-11e4-9089-e8e8daec308c.gif)


About ttyrec
------------

see [http://0xcc.net/ttyrec/](http://0xcc.net/ttyrec/).


Installation
------------

    go get github.com/lengshuiyulangcn/ttyrec2gif


Usage
-----

    ttyrec2gif -in <input file> -out <output file> -s <speed> -row <rows> -col <columns>

* `in`: ttyrec file (default: `"ttyrecord"`)
* `out`: output animated GIF file name (default: `"tty.gif"`)
* `s`: play speed (default: `1.0`)
* `row`: terminal height (default: `24`)
* `col`: terminal width (default: `80`)
* `noloop`: play only once (default: `false`)


Demo
====

256color, [peco](https://github.com/peco/peco)

![](https://cloud.githubusercontent.com/assets/80381/3735011/637e1d8a-171f-11e4-9d1b-4b5da5456378.gif)


sl

![](https://cloud.githubusercontent.com/assets/80381/3745301/b8836604-17a7-11e4-85d7-0441d875e486.gif)


Note
----

- [j4k.co/terminal](http://godoc.org/j4k.co/terminal)
- [Anonymous Pro](http://www.marksimonson.com/fonts/view/anonymous-pro) font
