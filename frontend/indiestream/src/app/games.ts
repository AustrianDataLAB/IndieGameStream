export interface Games {
  Message: string,
  Status: string,
  Games: Game[];
}

export interface Game {
  id: string,
  title: string,
  status: string,
  url: string
}
