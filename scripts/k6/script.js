import http from "k6/http";
import { check, sleep } from "k6";
import * as chance from "./chance.js";
import { randomLanguage, randomStack } from "./languages.js";

const c = chance.Chance();

function makePerson() {
  const person = {
    apelido: `${c.first()}_${c.word({ length: 10 })}`.toLowerCase(),
    nome: c.name(),
    nascimento: c.birthday().toISOString().substring(0, 10),
    stack: randomStack(c),
  };

  return person;
}

function makeTerm() {
  const factories = [
    () => c.first().toLowerCase(),
    () => c.name(),
    () => randomLanguage(c),
  ];

  const randomFactory =
    factories[c.integer({ min: 0, max: factories.length - 1 })];

  return randomFactory();
}

export const options = {
  stages: [{ duration: "5m", target: 1000 }],
};

export default function () {
  const baseUrl = "http://localhost:9999";

  let res1 = http.post(`${baseUrl}/pessoas`, JSON.stringify(makePerson()), {
    verb: "post",
    tags: { name: "Criação" },
    headers: {
      "Content-Type": "application/json",
    },
  });
  check(res1, {
    "status is 201": (r) => r.status === 201,
  });

  sleep(1);

  let res2 = http.get(`${baseUrl}${res1.headers["Location"]}`, {
    verb: "get",
    tags: { name: "Consulta" },
    headers: {
      "Content-Type": "application/json",
    },
  });
  check(res2, {
    "status is 200": (r) => r.status === 200,
  });

  sleep(1);

  let randomTerm = makeTerm();

  let res3 = http.get(`${baseUrl}/pessoas?t=${randomTerm}`, {
    verb: "get",
    tags: { name: "Busca" },
    headers: {
      "Content-Type": "application/json",
    },
  });
  check(res3, {
    "status is 200": (r) => r.status === 200,
  });
}
