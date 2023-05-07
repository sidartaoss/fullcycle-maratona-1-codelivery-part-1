# Maratona Full Cycle - Codelivery - Part I

O projeto consiste em:

- Um sistema de monitoramento de veículos de entrega em tempo real.

Requisitos:

- Uma transportadora quer fazer o agendamento de suas entregas;
- Ela também quer ter o feedback instantâneo de quando a entrega é realizada;
- Caso haja necessidade de acompanhar a entrega com mais detalhes, o sistema deverá informar, em tempo real, a localização do motorista no mapa.

#### Que problemas de negócio o projeto poderia resolver?
- O projeto pode ser adaptado para casos de uso onde é necessário rastrear e monitorar carros, caminhões, frotas e remessas em tempo real, como na logística e na indústria automotiva.

Dinâmica do sistema:

1. A aplicação Order (React/Nest.js) é responsável pelas ordens de serviço (ou pedidos) e vai conter a tela de agendamento de pedidos de entrega. A criação de uma nova ordem de serviço começa o processo para que o motorista entregue a mercadoria;

2. A aplicação Driver (Go) é responsável por gerenciar o contexto limitado de motoristas. Neste caso, sua responsabilidade consiste em disponibilizar os endpoints de consulta;

3. Para a criação de uma nova ordem de serviço, a aplicação Order obtém de Driver os dados dos motoristas. Neste caso, REST é uma opção pertinente, porque a comunicação deve ser a mais simples possível. Então, ela acontece de forma direta e síncrona com uma requisição HTTP;

4. Após criar a nova ordem de serviço, Order notifica a aplicação Mapping (Nest.js/React) via RabbitMQ de que o motorista deve iniciar a entrega. Mapping é a aplicação que vai exibir no mapa a posição do motorista em tempo real. A aplicação Simulator (Go) também é notificada sobre o início da entrega e começa a enviar para a aplicação Mapping as posições do veículo;

5. Ao finalizar a entrega, a aplicação Mapping notifica via RabbitMQ a aplicação Order de que o produto foi entregue e a aplicação altera o status da entrega de Pendente para Entregue.

## Tecnologias

#### Operate What You Build

- Nesta versão inicial, trabalhamos apenas com o backend. Posteriormente, serão adicionadas as tecnologias de frontend, integração contínua, deploy e observabilidade.

  - Backend
  
    - Golang
    - RabbitMQ

## Formatos de Comunicação

- REST
- Sistema de mensageria (RabbitMQ)

> Comunicação entre Serviços

- REST é um dos formatos de comunicação mais conhecidos, principalmente pelo uso em cenários aonde o cliente da requisição é o browser. No entanto, há situações em que o formato de comunicação entre serviços precisa ser outro. Por exemplo, quando há a necessidade da garantia de entrega das mensagens. Neste caso, um sistema de mensageria é uma opção bastante pertinente.

  - Quais são alguns dos sistemas de mensageria comuns no mercado?

    - RabbitMQ;
    - Apache Kafka (Sistema de Stream);
    - Amazon SQS;
    - ActiveMQ;
    - Redis.

> RabbitMQ

- RabbitMQ é um dos sistemas de mensageria mais tradicionais; é simples, funciona muito bem e vem sendo utilizado por sistemas gigantes. Costuma-se dizer que é como o HTML dos sistemas de fila, porque todo programador deveria conhecer, independente de ser a opção adotada na empresa.

  - Funcionamento Básico

    - O RabbitMQ trabalha com um Publisher e um Consumer;
    - O Publisher publica uma mensagem para ser consumida pelo Consumer;
    - A mensagem é publicada para uma fila, só que não é enviada diretamente para a fila - ela é enviada para uma Exchange;
    - E o que faz a Exchange? A Exchange contém regras e, baseado nessas regras, as mensagens são encaminhadas para as filas que foram configuradas.

  - Principais Tipos de Exchange

    - Direct: Com uma configuração, é possível escolher as filas para onde as mensagens vão cair;
    - Fanout: Todas as filas concentradas nessa Exchange vão receber a mensagem. A Exchange manda para todas as filas;
    - Topic: É uma Exchange similar à Direct, mas é mais flexível de configurar para qual fila enviar as mensagens.

  - Direct Exchange

    - O uso desse tipo de Exchange é bastante comum, sendo também a Exchange mais utilizada neste projeto. 
    - Qual é a idéia principal?
      - A partir de um Publisher (ou Producer), há vários Consumers que podem estar interessados em receber uma mensagem. Mas, cada Consumer vai consumir apenas a mensagem que fizer mais sentido para ele. Ou seja, ele fica escutando filas específicas. No entanto, não se trata de um Fanout: ao enviar a mensagem para a Exchange, a intenção não é replicar para todas as filas. Então, trabalha-se com uma regra. Essa regra vai permitir que a mensagem seja enviada de acordo com o que foi configurado.
    - E como acontece a configuração?
      - Através de um Bind. Bind é o relacionamento entre a fila e a Exchange. Para criar esse relacionamento, é configurando um Routing Key, ou seja, uma chave de roteamento;
      - Então, para cada fila, é configurado uma Routing Key e, no momento de enviar uma mensagem, ela é enviada para uma Exchange com a Routing Key X;
      - E, na seqüência, a Exchange encaminha para a fila que estiver vinculada à Routing Key X;
      - Assim, a mensagem é sempre enviada para a mesma Exchange, mas a Routing Key é diferente, ou seja, de acordo com a Routing Key, a Exchange pode mandar para uma fila diferente.

### Simulator

Optou-se pela linguagem Go, neste caso, porque resolve de forma muito simples a complexidade técnica da aplicação; o Go permite trabalhar muito facilmente com multithreading e é bastante performático.

Caso a aplicação Simulator fosse single-threaded, por exemplo, seria necessário aguardar o processamento de um veículo entregando para, então, começar a processar a entrega de outro veículo. Como a aplicação é multi-threaded, ela consegue gerenciar todos os veículos ao mesmo tempo com quantas threads forem necessárias: não importa quantas entregas forem criadas ao mesmo tempo.

Interessante notar a utilização do padrão Worker Pool na aplicação.

- Worker pool

  - Worker pool, também conhecido como Thread Pool, é um padrão de concorrência comumente utilizado no Go, no qual um número fixo de workers roda em paralelo para processar um número de tarefas que estão aguardando em uma fila. O Go faz uso de goroutines e channels para construir esse padrão. Normalmente, os workers são definidos por uma goroutine que fica aguardando até que obtenha dados através de um channel que é responsável por coordenar os workers e a tarefa na fila. Fonte: https://blog.devgenius.io/golang-concurrency-worker-pool-2aff9cbc6255

  - Por que utilizar o padrão Worker Pool?

    - Porque uma máquina não conta com recursos ilimitados. O tamanho mínimo de um objeto goroutine é de 2KB; quando muitas goroutines são geradas, a máquina fica rapidamente sem memória e a CPU continua processando a tarefa até atingir o limite. Ao utilizar um pool limitado de workers e manter a tarefa na fila, a possibilidade de estourar CPU e memória é reduzida, pois a tarefa aguarda na fila até que o worker puxe a tarefa. Fonte: https://medium.com/code-chasm/go-concurrency-pattern-worker-pool-a437117025b1

### Execução

#### RabbitMQ

1. Criar uma fila positions - responsável por notificar o Simulator sobre uma nova ordem de serviço;

2. Criar uma fila micro-mapping/new-position - responsável por consumir as posições enviadas pelo Simulator;

3. Adicionar um Bind à fila micro-mapping/new-position - From exchange: amq.direct; Routing key: mapping.new-position;

4. Publicar uma nova mensagem na fila positions.

![Queue positions - Publish message](./images/positions-publish-message.png)

#### Simulator

5. Executar o comando: `docker-compose up -d`;

6. Executar o comando: `docker-compose exec goapp_simulator bash`;

7. Executar o comando: `go run simulator.go`;

8. Acompanhar os logs de mensagens publicadas para o RabbitMQ com os dados de latitude e longitude. Ao enviar latitude 0 e longitude 0, sinaliza que a última posição foi enviada;

![Simulator - messages sent to RabbitMQ](./images/simulator-messages-sent-to-rabbitmq.png)

#### RabbitMQ

9. A fila micro-mapping/new-position permite acompanhar o envio de mensagens. Ir em Get message / Messages: 10 / clicar Get Message(s).

![Queue micro-mapping/new-position: Get Message(s)](./images/micro-mapping-new-position-get-messages.png)
