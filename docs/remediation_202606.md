# Correções da Análise Técnica - Junho/2026 (catnet)

Este documento sumariza a execução do plano de ação para sanar as vulnerabilidades e melhorias técnicas apontadas especificamente para a CLI `catnet`.

## Resumo das Modificações no `catnet`

1. **Testes de Integração Independentes (`C6`, `C7`, `C10`)**
   Substituímos a execução em-processo via `cli.Execute()` e a mutação de estado global (`os.Args`/`os.Stdout`) por chamadas a subprocessos isolados.
   - Criamos a função `TestMain` em `integration_test.go` que compila temporariamente o binário do CLI via `go build`.
   - Modificamos todos os testes para usar o pacote `os/exec`. Isso garante isolamento perfeito, permitindo até a simulação precisa de cancelamento por sinais de OS (`os.Interrupt`) sem desestabilizar ou matar a suíte de testes concorrente.
   - Foram eliminados os problemas de estado do `rootCmd` compartilhado entre testes.

2. **Isolamento de Flag (`C8`)**
   - Alteramos a arquitetura de flags para o comando `--format`.
   - Ao invés de defini-lo como `PersistentFlag` em `rootCmd` (o que causava herança indesejada onde formatos diferentes eram suportados), passamos a registrá-lo unicamente para os subcomandos específicos. O `exportCmd` agora possui sua flag dedicada sem conflito.

3. **Injeção de Dependência e Testabilidade (`C9`/`M7`)**
   - Refatoramos a estrutura `HumanOutput` (`output/human.go`) para utilizar atributos customizados do tipo `io.Writer` em vez de hardcodar `os.Stdout` e `os.Stderr`.
   - Criamos o pacote de testes em `human_test.go` para atestar com sucesso a formatação, e supressão de logs em modo silencioso (`--quiet`), verificando as injeções com buffers da memória.
