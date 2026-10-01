package main

import (
	"testing"
	"github.com/go-rod/rod"
	"time"
	"strings"
)


func TestAdicionarLivroAColecao(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage()
	fazerLogin(page)
	limparColecao(t, page)

	page.MustNavigate("https://demoqa.com/books")
	page.Timeout(15 * time.Second).MustElementR("a", "Speaking JavaScript").MustClick()

	capturarAlertas(page)
	botao := page.Timeout(15 * time.Second).MustElementR("button", "Add To Your Collection")
	botao.MustEval(`() => this.click()`)

	mensagem := esperarAlerta(page)
	if !strings.Contains(mensagem, "added") {
		t.Fatalf("esperava confirmação de livro adicionado, mas recebeu: %q", mensagem)
	}

	page.MustNavigate("https://demoqa.com/profile")
	page.Timeout(15 * time.Second).MustElementR("a", "Speaking JavaScript")
}

func TestNaoDuplicarLivroNaColecao(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage()
	fazerLogin(page)
	limparColecao(t, page)

	page.MustNavigate("https://demoqa.com/books")
	page.Timeout(15 * time.Second).MustElementR("a", "Speaking JavaScript").MustClick()

	capturarAlertas(page)
	botao := page.Timeout(15 * time.Second).MustElementR("button", "Add To Your Collection")

	botao.MustEval(`() => this.click()`)
	if mensagem := esperarAlerta(page); !strings.Contains(mensagem, "added") {
		t.Fatalf("esperava que a primeira adição fosse aceita, mas recebeu: %q", mensagem)
	}

	botao.MustEval(`() => this.click()`)
	mensagemDuplicidade := esperarAlerta(page)
	t.Logf("mensagem de duplicidade: %q", mensagemDuplicidade)

	if !strings.Contains(mensagemDuplicidade, "already present") {
		t.Fatalf("esperava mensagem de duplicidade, mas recebeu: %q", mensagemDuplicidade)
	}

	page.MustNavigate("https://demoqa.com/profile")
	page.Timeout(15 * time.Second).MustElementX(`//a[contains(text(),"Speaking JavaScript")]`)
	livros := page.MustElementsX(`//a[contains(text(),"Speaking JavaScript")]`)

	if len(livros) != 1 {
		t.Fatalf("esperava exatamente 1 livro 'Speaking JavaScript' na coleção, mas encontrou: %d", len(livros))
	}
}

func TestAcessarPerfilSemLogin(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/profile")

	mensagem := page.Timeout(15 * time.Second).
		MustElementR("*", "Currently you are not logged into the Book Store application").MustText()

	if !strings.Contains(mensagem, "Currently you are not logged") {
		t.Fatalf("esperava mensagem de acesso negado, mas recebeu: %s", mensagem)
	}
}

func TestRemoverLivroDaColecao(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage()
	fazerLogin(page)
	limparColecao(t, page)
	adicionarLivrosViaAPI(t, page, "9781449325862", "9781449365035")

	page.MustNavigate("https://demoqa.com/profile")
	page.Timeout(15 * time.Second).MustElementR("a", "Git Pocket Guide")
	page.Timeout(15 * time.Second).MustElementR("a", "Speaking JavaScript")

	capturarAlertas(page)
	page.MustElement("#delete-record-9781449365035").MustEval(`() => this.click()`)
	page.Timeout(10 * time.Second).MustElement("#closeSmallModal-ok").MustEval(`() => this.click()`)

	if mensagem := esperarAlerta(page); !strings.Contains(mensagem, "deleted") {
		t.Fatalf("esperava confirmação de livro excluído, mas recebeu: %q", mensagem)
	}

	page.MustNavigate("https://demoqa.com/profile")
	page.Timeout(15 * time.Second).MustElementR("a", "Git Pocket Guide")

	if page.MustHasR("a", "Speaking JavaScript") {
		t.Fatal("esperava que 'Speaking JavaScript' tivesse sido removido, mas ele continua na coleção após recarregar o perfil")
	}

	if quantidade := contarLivros(page); quantidade != 1 {
		t.Fatalf("esperava 1 livro restante na coleção, mas encontrou: %d", quantidade)
	}
}

func TestBotaoAdicionarOcultoSemLogin(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/books")
	page.Timeout(15 * time.Second).MustElementR("a", "Speaking JavaScript").MustClick()
	page.Timeout(15 * time.Second).MustElementR("*", "9781449365035")

	if page.MustHasR("button", "Add To Your Collection") {
		t.Fatal("o botão 'Add To Your Collection' está disponível para um visitante sem login")
	}
}