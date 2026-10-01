package main

import (
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
)

func TestBuscarLivroExistente(t *testing.T) {

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/books")

	campoBusca := page.Timeout(15 * time.Second).MustElement("input[placeholder='Type to search']")
	page.Timeout(15 * time.Second).MustElementR("a", "Git Pocket Guide")

	antes := contarLivros(page)
	if antes <= 1 {
		t.Fatalf("esperava vários livros listados antes da busca, mas o seletor encontrou: %d", antes)
	}

	campoBusca.MustInput("Git")

	depois := esperarQuantidadeDeLivros(page, 1)
	if depois != 1 {
		t.Fatalf("esperava 1 livro após buscar 'Git', mas encontrou: %d (antes da busca: %d)", depois, antes)
	}

	page.Timeout(5 * time.Second).MustElementR("a", "Git Pocket Guide")
}

func TestBuscarTermoInexistente(t *testing.T) {

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/books")

	campoBusca := page.Timeout(15 * time.Second).MustElement("input[placeholder='Type to search']")

	campoBusca.MustInput("xyznaoexiste123")


	paginacao := page.Timeout(10 * time.Second).MustElementR("*", "Page 1 of 0").MustText()

	if !strings.Contains(paginacao, "Page 1 of 0") {
		t.Fatalf("esperava paginação 'Page 1 of 0', mas recebeu: %s", paginacao)
	}
}

func TestBuscarIgnoraMaiusculasEMinusculas(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/books")

	campoBusca := page.Timeout(15 * time.Second).MustElement("input[placeholder='Type to search']")
	page.Timeout(15 * time.Second).MustElementR("a", "Git Pocket Guide")

	campoBusca.MustInput("gIT")

	if quantidade := esperarQuantidadeDeLivros(page, 1); quantidade != 1 {
		t.Fatalf("esperava 1 livro ao buscar 'gIT', como ao buscar 'Git', mas encontrou: %d", quantidade)
	}

	page.Timeout(5 * time.Second).MustElementR("a", "Git Pocket Guide")
}

func TestBuscarComCaracteresEspeciais(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/books")

	campoBusca := page.Timeout(15 * time.Second).MustElement("input[placeholder='Type to search']")
	page.Timeout(15 * time.Second).MustElementR("a", "Git Pocket Guide")

	for _, termo := range []string{"(", "[", "*", `\`} {
		campoBusca.MustSelectAllText().MustInput(termo)
		if _, err := page.Timeout(10 * time.Second).ElementR("*", "Page 1 of 0"); err != nil {
			t.Fatalf("ao buscar %q, esperava a tabela funcionando e vazia ('Page 1 of 0'), mas: %v", termo, err)
		}

		campoBusca.MustSelectAllText().MustInput("Git")
		if quantidade := esperarQuantidadeDeLivros(page, 1); quantidade != 1 {
			t.Fatalf("depois de buscar %q, a busca por 'Git' deveria retornar 1 livro, mas retornou: %d", termo, quantidade)
		}
	}
}