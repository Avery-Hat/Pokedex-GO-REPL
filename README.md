# Pokedex CLI: Map, Explore, Catch & Inspect

This CLI lets you browse Pokémon **location areas**, explore which Pokémon appear there, and build a personal **Pokédex** using the public [PokeAPI](https://pokeapi.co/).

## Run

```bash
go run .
```

You’ll see a prompt:

```
Pokedex >
```

Type commands and press **Enter**.

## Commands

### `map`: list location areas (20 at a time)

Shows the next page of 20 location areas each time you run it.

```text
Pokedex > map
canalave-city-area
eterna-city-area
pastoria-city-area
sunyshore-city-area
sinnoh-pokemon-league-area
...
```

Run `map` again to see the next 20:

```text
Pokedex > map
mt-coronet-1f-route-216
mt-coronet-1f-route-211
great-marsh-area-1
...
```

### `mapb`: go back one page

If you’ve paged forward with `map`, use `mapb` to return to the previous 20.
If you’re already on the first page, it prints:

```
you're on the first page
```

### `explore <area-name>`: list Pokémon in an area

Pass any area name you saw from `map`/`mapb` (use it **exactly** as shown).

```text
Pokedex > explore pastoria-city-area
Exploring pastoria-city-area...
Found Pokemon:
 - tentacool
 - tentacruel
 - magikarp
 - gyarados
 - remoraid
 - octillery
 - wingull
 - pelipper
 - shellos
 - gastrodon
```

### `catch <pokemon-name>`: try to catch a Pokémon

Uses the Pokémon endpoint and a catch roll (harder if base experience is high).

```text
Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu escaped!

Pokedex > catch pikachu
Throwing a Pokeball at pikachu...
pikachu was caught!
You may now inspect it with the inspect command.
```

Tip: multiword names should use dashes (e.g., `mr-mime`). The CLI accepts `mr mime` too and converts it.

### `inspect <pokemon-name>`: show details for a caught Pokémon

Prints name, height, weight, stats, and types **from your local Pokédex** (no API call).

```text
Pokedex > inspect pidgey
you have not caught that pokemon

Pokedex > catch pidgey
Throwing a Pokeball at pidgey...
pidgey was caught!
You may now inspect it with the inspect command.

Pokedex > inspect pidgey
Name: pidgey
Height: 3
Weight: 18
Stats:
  - hp: 40
  - attack: 45
  - defense: 40
  - special-attack: 35
  - special-defense: 35
  - speed: 56
Types:
  - normal
  - flying
```

### `pokedex`: list everything you’ve caught

```text
Pokedex > pokedex
Your Pokedex:
 - caterpie
 - pidgey
```

## Tips

* **Case-insensitive input:** Commands and names are normalized to lowercase.
* **Caching:** The first fetch of a page/area/Pokémon may pause briefly (network). Repeating the same request should be **instant** thanks to the cache.
* **Help & exit:**

  * `help`: shows available commands
  * `exit`: quit the CLI

## Optional: log your session

```bash
go run . | tee repl.log
```

Search later:

```bash
grep -i "forest" repl.log
```

---
