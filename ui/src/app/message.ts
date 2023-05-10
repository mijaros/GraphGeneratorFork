export enum Severity {
  ERROR = 'danger',
  WARNING = 'warning',
  SUCCESS = 'success'
}

export class Message {
  constructor(public message: string, public severity: Severity) {
  }

  get messageClass(): string {
    switch (this.severity) {
      case Severity.ERROR:
        return "bg-danger"
      case Severity.WARNING:
        return "bg-warning"
      case Severity.SUCCESS:
        return "bg-success"
    }
  }

}
