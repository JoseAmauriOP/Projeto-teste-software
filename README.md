# Testes automatizados end-to-end com Go e Rod

Suíte de **16 testes end-to-end** escrita em **Go** que controla um navegador real com a biblioteca [Rod](https://github.com/go-rod/rod) para testar a loja de livros do [DemoQA](https://demoqa.com/books), um site público feito para praticar automação de testes.

Projeto da disciplina de Teste de Software (UFC), usado na apresentação de ferramentas de teste.

## O que é testado

| Arquivo | Área | Testes |
|---|---|---|
| `login_test.go` | Login e logout | login válido, credenciais inválidas, campos vazios, senha com maiúsculas trocadas, logout encerra a sessão |
| `busca_test.go` | Busca de livros | termo existente, termo inexistente, maiúsculas e minúsculas, caracteres especiais (`(`, `[`, `*`, `\`) |
| `colecao_test.go` | Coleção do usuário | adicionar livro, não duplicar livro, remover livro, perfil sem login, botão de adicionar oculto sem login |
| `Pagina_test.go` | Página do livro | carregamento da loja, ISBN inexistente |
| `auxiliares_test.go` | Funções de apoio | login, leitura da sessão, chamadas à API, espera por alertas e elementos |

Alguns testes foram escritos para **investigar defeitos** do site, e não só confirmar o que funciona. Por exemplo, o teste de ISBN inexistente verifica se o sistema avisa que o livro não existe, ou se mostra uma ficha vazia.

## Decisões de projeto

- **Estado limpo antes de cada teste**: os testes de coleção apagam a coleção do usuário pela API REST do site antes de começar, para não depender do resultado de outro teste.
- **Preparação pela API, verificação pela tela**: quando o teste precisa de livros já na coleção, eles são adicionados direto pela API (`POST /BookStore/v1/Books`), que é mais rápido e estável que clicar pela interface.
- **Esperas explícitas**: em vez de pausas fixas, os testes esperam o elemento ou a condição aparecer, com tempo limite.
- **Alertas do navegador**: o `window.alert` é substituído por uma função que guarda a mensagem, para o teste conseguir ler a confirmação.
- **Mensagens de falha descritivas**: cada falha diz o que era esperado e o que foi encontrado.

## Tecnologias

- Go e o pacote `testing`
- [Rod](https://github.com/go-rod/rod): automação do Chromium pelo DevTools Protocol
- [gotestsum](https://github.com/gotestyourself/gotestsum): saída de testes mais legível
- `net/http` para as chamadas à API

## Como rodar

Na primeira execução, o Rod baixa o Chromium automaticamente. É preciso ter internet, porque os testes acessam o site do DemoQA.

```bash
git clone https://github.com/JoseAmauriOP/Projeto-teste-software.git
cd Projeto-teste-software

go test -v ./...
```

Com saída formatada pelo gotestsum:

```bash
go tool gotestsum --format testdox
```

Para rodar só um grupo de testes:

```bash
go test -v -run TestLogin ./...
```

Os testes de login e coleção usam uma conta de teste do próprio DemoQA. Se ela deixar de existir, crie outra em `https://demoqa.com/register` e troque as constantes `usuarioValido` e `senhaValida` em `auxiliares_test.go`.

## O que aprendi

- Escrever testes end-to-end que controlam um navegador de verdade
- Escolher seletores (CSS, XPath, texto) que não quebram à toa
- Lidar com páginas assíncronas usando esperas com tempo limite
- Combinar teste de interface com chamadas de API para preparar o cenário
- Deixar cada teste independente dos outros
- Escrever testes que procuram defeitos, e não só confirmam o caminho feliz

---

Desenvolvido por [José Amauri Patrício](https://www.linkedin.com/in/amauri-patricio-5a121435b)