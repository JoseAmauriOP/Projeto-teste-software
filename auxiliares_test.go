package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
)

const (
	usuarioValido = "contateste"
	senhaValida   = "Teste@123"
)

func fazerLogin(page *rod.Page) {
	page.MustNavigate("https://demoqa.com/login")
	page.Timeout(15 * time.Second).MustElement("#userName").MustInput(usuarioValido)
	page.MustElement("#password").MustInput(senhaValida)
	page.MustElement("#login").MustClick()
	page.Timeout(15 * time.Second).MustElementR("button", "Logout")
}

func lerSessao(t *testing.T, page *rod.Page) (string, string) {
	t.Helper()

	cookies := map[string]string{}
	for _, c := range page.MustCookies() {
		cookies[c.Name] = c.Value
	}

	userID, token := cookies["userID"], cookies["token"]
	if userID == "" || token == "" {
		t.Fatalf("cookies de sessão ausentes: userID=%q, token presente=%v", userID, token != "")
	}
	return userID, token
}

func chamarAPI(t *testing.T, token, metodo, url string, corpo []byte, statusEsperado int) {
	t.Helper()

	req, err := http.NewRequest(metodo, url, bytes.NewReader(corpo))
	if err != nil {
		t.Fatalf("falha ao montar a requisição %s %s: %v", metodo, url, err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("falha ao chamar %s %s: %v", metodo, url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != statusEsperado {
		resposta, _ := io.ReadAll(resp.Body)
		t.Fatalf("esperava status %d em %s %s, mas recebeu: %d %s", statusEsperado, metodo, url, resp.StatusCode, resposta)
	}
}

func limparColecao(t *testing.T, page *rod.Page) {
	t.Helper()
	userID, token := lerSessao(t, page)
	chamarAPI(t, token, http.MethodDelete, "https://demoqa.com/BookStore/v1/Books?UserId="+userID, nil, http.StatusNoContent)
}

func adicionarLivrosViaAPI(t *testing.T, page *rod.Page, isbns ...string) {
	t.Helper()
	userID, token := lerSessao(t, page)

	colecao := []map[string]string{}
	for _, isbn := range isbns {
		colecao = append(colecao, map[string]string{"isbn": isbn})
	}

	corpo, err := json.Marshal(map[string]any{"userId": userID, "collectionOfIsbns": colecao})
	if err != nil {
		t.Fatalf("falha ao montar o corpo da requisição: %v", err)
	}
	chamarAPI(t, token, http.MethodPost, "https://demoqa.com/BookStore/v1/Books", corpo, http.StatusCreated)
}

func capturarAlertas(page *rod.Page) {
	page.MustEval(`() => { window.__ultimoAlerta = ""; window.alert = (msg) => { window.__ultimoAlerta = msg } }`)
}

func esperarAlerta(page *rod.Page) string {
	limite := time.Now().Add(10 * time.Second)
	for time.Now().Before(limite) {
		msg := page.MustEval(`() => window.__ultimoAlerta`).String()
		if msg != "" {
			page.MustEval(`() => { window.__ultimoAlerta = "" }`)
			return msg
		}
		time.Sleep(300 * time.Millisecond)
	}
	return ""
}

func contarLivros(page *rod.Page) int {
	return len(page.MustElementsX(`//table/tbody/tr[.//a]`))
}

func esperarQuantidadeDeLivros(page *rod.Page, esperado int) int {
	limite := time.Now().Add(10 * time.Second)
	quantidade := -1
	for time.Now().Before(limite) {
		quantidade = contarLivros(page)
		if quantidade == esperado {
			return quantidade
		}
		time.Sleep(300 * time.Millisecond)
	}
	return quantidade
}

func esperarCampoInvalido(page *rod.Page, seletor string) (bool, string) {
	campo := page.MustElement(seletor)
	limite := time.Now().Add(10 * time.Second)
	classe := ""
	for time.Now().Before(limite) {
		if c := campo.MustAttribute("class"); c != nil {
			classe = *c
		}
		if strings.Contains(classe, "is-invalid") {
			return true, classe
		}
		time.Sleep(300 * time.Millisecond)
	}
	return false, classe
}

func lerCampoDoLivro(page *rod.Page, campo string) (bool, string) {
	linha := "#" + campo + "-wrapper"
	if !page.MustHas(linha) {
		return false, ""
	}

	existe, valor, err := page.Has(linha + " #userName-value")
	if err != nil || !existe {
		return true, ""
	}
	return true, strings.TrimSpace(valor.MustText())
}

func esperarCampoDoLivro(page *rod.Page, campo, esperado string) (bool, string) {
	limite := time.Now().Add(10 * time.Second)
	existe, valor := false, ""
	for time.Now().Before(limite) {
		if existe, valor = lerCampoDoLivro(page, campo); valor == esperado {
			return existe, valor
		}
		time.Sleep(300 * time.Millisecond)
	}
	return existe, valor
}