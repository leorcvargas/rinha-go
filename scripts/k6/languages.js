const LANGUAGES = [
  "C#",
  "Java",
  "JavaScript",
  "Python",
  "Ruby",
  "Go",
  "Rust",
  "C",
  "C++",
  "PHP",
  "TypeScript",
  "Kotlin",
  "Swift",
  "Objective-C",
  "Scala",
  "R",
  "Dart",
  "Elixir",
  "Clojure",
  "Haskell",
  "Julia",
  "Perl",
  "Lua",
  "Erlang",
  "F#",
  "Groovy",
  "Haxe",
  "CoffeeScript",
  "OCaml",
  "Scheme",
  "Visual Basic",
  "Assembly",
  "PowerShell",
  "Delphi",
];

const LANGUAGES_LENGTH = LANGUAGES.length;

export function randomLanguage(chance) {
  return LANGUAGES[chance.integer({ min: 0, max: LANGUAGES_LENGTH - 1 })];
}

export function randomStack(chance) {
  const stackSize = chance.integer({ min: 0, max: 5 });

  if (stackSize === 0) {
    return null;
  }

  const stack = [];
  for (let i = 0; i < stackSize; i++) {
    stack.push(randomLanguage(chance));
  }

  return stack;
}
