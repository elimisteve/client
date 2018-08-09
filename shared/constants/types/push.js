// @flow
import * as I from 'immutable'
import * as ChatTypes from './chat2'
import * as RPCChatTypes from './rpc-chat-gen'
export type TokenType = 'apple' | 'appledev' | 'androidplay'

export type PushNotification =
  | {
      type: 'chat.readmessage',
      badges: number,
    }
  | {
      type: 'chat.newmessageSilent_2',
    }
  | {
      type: 'chat.newmessage',
      conversationIDKey: ChatTypes.ConversationIDKey,
      userInteraction: boolean,
      membersType: RPCChatTypes.ConversationMembersType,
      unboxPayload: string,
    }
  | {
      type: 'follow',
      userInteraction: boolean,
      username: string,
    }

export type _State = {
  token: string,
  showPushPrompt: boolean,
  hasPermissions: boolean,
}

export type State = I.RecordOf<_State>
