const payload = {
  nome: "Teste Frontend",
  tipo: "CREDITO",
  bandeira: "Visa",
  taxa: "2.5",
  taxa_fixa: "0.5",
  d_mais: 30,
  ordem_exibicao: 0,
  ativo: true
};

console.log(JSON.stringify(payload, null, 2));
