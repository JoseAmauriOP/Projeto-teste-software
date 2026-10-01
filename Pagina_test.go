package main

import (
	"testing"
	"time"
	"github.com/go-rod/rod"
)

func TestAbrirPaginaBookStore(t *testing.T) {
	// 1. Conecta ao navegador (baixa Chromium automaticamente se necessário)
	browser := rod.New().MustConnect()

	defer browser.MustClose()

	// 2. Abre a página do Book Store no demoqa
	page := browser.MustPage("https://demoqa.com/books")

	// 3. Espera até 15s pelo campo de busca aparecer
	campoBusca := page.Timeout(15 * time.Second).MustElement("input[placeholder='Type to search']")

	// 4. Valida que o campo realmente está visível na tela
	if !campoBusca.MustVisible() {
		t.Fatal("campo de busca não está visível - página pode não ter carregado corretamente")
	}

}

func TestAbrirLivroComIsbnInexistente(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/books?search=9781449365035")
	existe, titulo := esperarCampoDoLivro(page, "title", "Speaking JavaScript")
	if !existe {
		t.Fatal("calibração: a linha do título (#title-wrapper) não foi encontrada na ficha de um ISBN válido")
	}
	if titulo != "Speaking JavaScript" {
		t.Fatalf("calibração: esperava ler 'Speaking JavaScript' no campo Title de um ISBN válido, mas leu: %q", titulo)
	}

	page.MustNavigate("https://demoqa.com/books?search=9780000000000")
	page.Timeout(15 * time.Second).MustElementR("button", "Back To Book Store")

	limite := time.Now().Add(10 * time.Second)
	for time.Now().Before(limite) {
		if existe, _ := lerCampoDoLivro(page, "title"); !existe {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}

	if _, titulo := lerCampoDoLivro(page, "title"); titulo != "" {
		t.Fatalf("para o ISBN inexistente 9780000000000, o sistema exibiu o título %q", titulo)
	}
	t.Fatal("para o ISBN inexistente 9780000000000, o sistema exibiu a ficha do livro com os campos vazios, sem informar que o livro não existe")
}