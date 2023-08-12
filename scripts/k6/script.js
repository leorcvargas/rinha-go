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
  stages: [
    { duration: "10m", target: 1000 },
  ],
};

export default function () {
  const baseUrl = "http://localhost:9999";

  let res1 = http.post(`${baseUrl}/pessoas`, JSON.stringify(makePerson()), {
    verb: "post",
    headers: {
      "Content-Type": "application/json",
    },
  });
  check(res1, {
    "status is 201": (r) => r.status === 201,
  });

  // // sleep(1);

  // const location = res1.headers["Location"];

  // let res2 = http.get(`${baseUrl}${location}`, {
  //   verb: "get",
  //   tags: { name: "get_person" },
  // });

  // check(res2, {
  //   "status is 200": (r) => r.status === 200,
  // });

  // // sleep(1);

  // let res3 = http.get(encodeURI(`${baseUrl}/pessoas?t=${makeTerm()}`), {
  //   verb: "get",
  //   tags: { name: "search_person" },
  // });

  // check(res3, {
  //   "status is 200": (r) => r.status === 200,
  // });
}
