package main

import  (
	"testing"
	"time"
	"github.com/go-rod/rod"
	"strings"
)

func TestLoginComSucesso(t *testing.T) {

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/login")

	page.Timeout(15 * time.Second).MustElement("#userName").MustInput("contateste")
	page.MustElement("#password").MustInput("Teste@123")

	page.MustElement("#login").MustClick()

	logout := page.Timeout(15 * time.Second).MustElementR("button", "Logout").MustText()

	if !strings.Contains(logout, "Logout") {
		t.Fatalf("esperava estar logado (botão 'Logout' visível), mas recebeu: %s", logout)
	}
}

func TestLoginComCredenciaisInvalidas(t *testing.T) {

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/login")

	page.Timeout(15 * time.Second).MustElement("#userName").MustInput("usuario_invalido")
	page.MustElement("#password").MustInput("SenhaErrada@123")

	page.MustElement("#login").MustClick()

	err := page.Timeout(15 * time.Second).MustElementR("*", "Invalid username or password!").MustText()

	if !strings.Contains(err, "Invalid username or password!") {
		t.Fatalf("esperava mensagem de erro, mas recebeu: %s", err)
	}
}

func TestLoginComCamposVazios(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("https://demoqa.com/login")

	page.Timeout(15 * time.Second).MustElement("#login").MustClick()

	for _, campo := range []string{"#userName", "#password"} {
		if invalido, classe := esperarCampoInvalido(page, campo); !invalido {
			t.Fatalf("esperava o campo %s marcado como inválido, mas suas classes são: %q", campo, classe)
		}
	}

	url := page.MustInfo().URL

	if !strings.Contains(url, "/login") {
		t.Fatalf("esperava continuar na página de login após campos vazios, mas foi para: %s", url)
	}
}

func TestLogoutEncerraSessao(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage()
	fazerLogin(page)

	page.MustElementR("button", "Logout").MustClick()

	page.MustNavigate("https://demoqa.com/profile")
	page.Timeout(15 * time.Second).MustElementR("*", "Currently you are not logged into the Book Store application")
}

func TestLoginComSenhaIncorreta(t *testing.T) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	senhaIncorreta := strings.ToLower(senhaValida)

	page := browser.MustPage("https://demoqa.com/login")

	page.Timeout(15 * time.Second).MustElement("#userName").MustInput(usuarioValido)
	page.MustElement("#password").MustInput(senhaIncorreta)
	page.MustElement("#login").MustClick()

	sessaoAberta := false
	page.Timeout(15 * time.Second).Race().
		ElementR("*", "Invalid username or password!").
		ElementR("button", "Logout").MustHandle(func(*rod.Element) { sessaoAberta = true }).
		MustDo()

	if sessaoAberta {
		t.Fatalf("o sistema aceitou a senha %q para o usuário %q", senhaIncorreta, usuarioValido)
	}

	if page.MustHasR("button", "Logout") {
		t.Fatal("o sistema exibiu a mensagem de erro, mas também abriu uma sessão com a senha incorreta")
	}
}