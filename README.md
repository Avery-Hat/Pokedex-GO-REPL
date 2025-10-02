# Pokedex CLI — Map & Explore Quickstart

This CLI lets you browse Pokémon **location areas** and explore which Pokémon appear in each area using the public [PokeAPI](https://pokeapi.co/).

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

### `map` — list location areas (20 at a time)

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

### `mapb` — go back one page

If you’ve paged forward with `map`, use `mapb` to return to the previous 20.
If you’re already on the first page, it prints:

```
you're on the first page
```

### `explore <area-name>` — list Pokémon in an area

Pass any area name you saw from `map`/`mapb`.

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

## Tips

* **Case-insensitive input:** Commands and area names are normalized to lowercase.
* **Caching:** The first fetch of a page/area may pause briefly (network). Re-running `mapb` or re-exploring the same area should be **instant** thanks to the cache.
* **Help & exit:**

  * `help` — shows available commands
  * `exit` — quit the CLI

## Optional: log your session

```bash
go run . | tee repl.log
```

You can later search it, e.g.:

```bash
grep -i "forest" repl.log
```

---
