import {Injectable} from '@angular/core';
import {Message, Severity} from "./message";

@Injectable({
  providedIn: 'root'
})
export class MessageService {
  messages: Message[] = []

  constructor() {
  }

  addMessage(message: string, severity: Severity) {
    let mess = new Message(message, severity)
    this.messages.push(mess)
  }

  addError(err: any) {
    console.log(err)
    let message = new Message(err.error.error.toString(), Severity.ERROR)
    this.messages.push(message)
  }

  removeMessage(message: Message) {
    this.messages = this.messages.filter((t) => t !== message)
  }
}
