/* eslint-disable */
import { grpc } from "@improbable-eng/grpc-web";
import { BrowserHeaders } from "browser-headers";
import _m0 from "protobufjs/minimal";

export const protobufPackage = "weshnet.push.v1";

export enum PushServiceTokenType {
  PushTokenUndefined = 0,
  /** PushTokenMQTT - PushTokenMQTT: Platform independent */
  PushTokenMQTT = 1,
  /** PushTokenApplePushNotificationService - PushTokenApplePushNotificationService: iOS, iPadOS, tvOS, macOS */
  PushTokenApplePushNotificationService = 2,
  /** PushTokenFirebaseCloudMessaging - PushTokenFirebaseCloudMessaging: Android with GMS, Chrome OS */
  PushTokenFirebaseCloudMessaging = 3,
  /** PushTokenWindowsPushNotificationService - PushTokenWindowsPushNotificationService: Windows, XBox */
  PushTokenWindowsPushNotificationService = 4,
  /** PushTokenHuaweiPushKit - PushTokenHuaweiPushKit: Huawei Android devices with AppGallery */
  PushTokenHuaweiPushKit = 5,
  /** PushTokenAmazonDeviceMessaging - PushTokenAmazonDeviceMessaging: Fire OS devices */
  PushTokenAmazonDeviceMessaging = 6,
  UNRECOGNIZED = -1,
}

export function pushServiceTokenTypeFromJSON(object: any): PushServiceTokenType {
  switch (object) {
    case 0:
    case "PushTokenUndefined":
      return PushServiceTokenType.PushTokenUndefined;
    case 1:
    case "PushTokenMQTT":
      return PushServiceTokenType.PushTokenMQTT;
    case 2:
    case "PushTokenApplePushNotificationService":
      return PushServiceTokenType.PushTokenApplePushNotificationService;
    case 3:
    case "PushTokenFirebaseCloudMessaging":
      return PushServiceTokenType.PushTokenFirebaseCloudMessaging;
    case 4:
    case "PushTokenWindowsPushNotificationService":
      return PushServiceTokenType.PushTokenWindowsPushNotificationService;
    case 5:
    case "PushTokenHuaweiPushKit":
      return PushServiceTokenType.PushTokenHuaweiPushKit;
    case 6:
    case "PushTokenAmazonDeviceMessaging":
      return PushServiceTokenType.PushTokenAmazonDeviceMessaging;
    case -1:
    case "UNRECOGNIZED":
    default:
      return PushServiceTokenType.UNRECOGNIZED;
  }
}

export function pushServiceTokenTypeToJSON(object: PushServiceTokenType): string {
  switch (object) {
    case PushServiceTokenType.PushTokenUndefined:
      return "PushTokenUndefined";
    case PushServiceTokenType.PushTokenMQTT:
      return "PushTokenMQTT";
    case PushServiceTokenType.PushTokenApplePushNotificationService:
      return "PushTokenApplePushNotificationService";
    case PushServiceTokenType.PushTokenFirebaseCloudMessaging:
      return "PushTokenFirebaseCloudMessaging";
    case PushServiceTokenType.PushTokenWindowsPushNotificationService:
      return "PushTokenWindowsPushNotificationService";
    case PushServiceTokenType.PushTokenHuaweiPushKit:
      return "PushTokenHuaweiPushKit";
    case PushServiceTokenType.PushTokenAmazonDeviceMessaging:
      return "PushTokenAmazonDeviceMessaging";
    case PushServiceTokenType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum PushServicePriority {
  PushPriorityUndefined = 0,
  PushPriorityLow = 1,
  PushPriorityNormal = 2,
  UNRECOGNIZED = -1,
}

export function pushServicePriorityFromJSON(object: any): PushServicePriority {
  switch (object) {
    case 0:
    case "PushPriorityUndefined":
      return PushServicePriority.PushPriorityUndefined;
    case 1:
    case "PushPriorityLow":
      return PushServicePriority.PushPriorityLow;
    case 2:
    case "PushPriorityNormal":
      return PushServicePriority.PushPriorityNormal;
    case -1:
    case "UNRECOGNIZED":
    default:
      return PushServicePriority.UNRECOGNIZED;
  }
}

export function pushServicePriorityToJSON(object: PushServicePriority): string {
  switch (object) {
    case PushServicePriority.PushPriorityUndefined:
      return "PushPriorityUndefined";
    case PushServicePriority.PushPriorityLow:
      return "PushPriorityLow";
    case PushServicePriority.PushPriorityNormal:
      return "PushPriorityNormal";
    case PushServicePriority.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface PushServiceServerInfo {
}

export interface PushServiceServerInfo_Request {
}

export interface PushServiceServerInfo_Reply {
  /**
   * public_key the public key used to encrypt data to the server (ie. a PushServiceReceiver),
   * to be used with NaCl's Seal and Open
   */
  publicKey: Uint8Array;
  /** supported_token_types list of token types supported by the server */
  supportedTokenTypes: PushServiceSupportedTokenType[];
}

export interface PushServiceSupportedTokenType {
  appBundleId: string;
  tokenType: PushServiceTokenType;
}

export interface PushServiceSend {
}

export interface PushServiceSend_Request {
  envelope: OutOfStoreMessageEnvelope | undefined;
  priority: PushServicePriority;
  receivers: PushServiceOpaqueReceiver[];
}

export interface PushServiceSend_Reply {
}

export interface OutOfStoreMessageEnvelope {
  nonce: Uint8Array;
  box: Uint8Array;
  groupReference: Uint8Array;
}

export interface OutOfStoreExposedData {
  nonce: Uint8Array;
  box: Uint8Array;
}

export interface PushServiceOpaqueReceiver {
  opaqueToken: Uint8Array;
  serviceAddr: string;
}

export interface DecryptedPush {
  accountId: string;
  accountName: string;
  conversationPublicKey: string;
  conversationDisplayName: string;
  memberPublicKey: string;
  memberDisplayName: string;
  pushType: DecryptedPush_PushType;
  payloadAttrsJson: string;
  deepLink: string;
  alreadyReceived: boolean;
  accountMuted: boolean;
  conversationMuted: boolean;
  hidePreview: boolean;
}

export enum DecryptedPush_PushType {
  Unknown = 0,
  Message = 1,
  GroupInvitation = 7,
  ConversationNameChanged = 8,
  MemberNameChanged = 9,
  MemberDetailsChanged = 11,
  UNRECOGNIZED = -1,
}

export function decryptedPush_PushTypeFromJSON(object: any): DecryptedPush_PushType {
  switch (object) {
    case 0:
    case "Unknown":
      return DecryptedPush_PushType.Unknown;
    case 1:
    case "Message":
      return DecryptedPush_PushType.Message;
    case 7:
    case "GroupInvitation":
      return DecryptedPush_PushType.GroupInvitation;
    case 8:
    case "ConversationNameChanged":
      return DecryptedPush_PushType.ConversationNameChanged;
    case 9:
    case "MemberNameChanged":
      return DecryptedPush_PushType.MemberNameChanged;
    case 11:
    case "MemberDetailsChanged":
      return DecryptedPush_PushType.MemberDetailsChanged;
    case -1:
    case "UNRECOGNIZED":
    default:
      return DecryptedPush_PushType.UNRECOGNIZED;
  }
}

export function decryptedPush_PushTypeToJSON(object: DecryptedPush_PushType): string {
  switch (object) {
    case DecryptedPush_PushType.Unknown:
      return "Unknown";
    case DecryptedPush_PushType.Message:
      return "Message";
    case DecryptedPush_PushType.GroupInvitation:
      return "GroupInvitation";
    case DecryptedPush_PushType.ConversationNameChanged:
      return "ConversationNameChanged";
    case DecryptedPush_PushType.MemberNameChanged:
      return "MemberNameChanged";
    case DecryptedPush_PushType.MemberDetailsChanged:
      return "MemberDetailsChanged";
    case DecryptedPush_PushType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface FormatedPush {
  pushType: DecryptedPush_PushType;
  title: string;
  subtitle: string;
  body: string;
  deepLink: string;
  muted: boolean;
  hidePreview: boolean;
  conversationIdentifier: string;
}

function createBasePushServiceServerInfo(): PushServiceServerInfo {
  return {};
}

export const PushServiceServerInfo = {
  encode(_: PushServiceServerInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PushServiceServerInfo {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePushServiceServerInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): PushServiceServerInfo {
    return {};
  },

  toJSON(_: PushServiceServerInfo): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<PushServiceServerInfo>, I>>(base?: I): PushServiceServerInfo {
    return PushServiceServerInfo.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<PushServiceServerInfo>, I>>(_: I): PushServiceServerInfo {
    const message = createBasePushServiceServerInfo();
    return message;
  },
};

function createBasePushServiceServerInfo_Request(): PushServiceServerInfo_Request {
  return {};
}

export const PushServiceServerInfo_Request = {
  encode(_: PushServiceServerInfo_Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PushServiceServerInfo_Request {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePushServiceServerInfo_Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): PushServiceServerInfo_Request {
    return {};
  },

  toJSON(_: PushServiceServerInfo_Request): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<PushServiceServerInfo_Request>, I>>(base?: I): PushServiceServerInfo_Request {
    return PushServiceServerInfo_Request.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<PushServiceServerInfo_Request>, I>>(_: I): PushServiceServerInfo_Request {
    const message = createBasePushServiceServerInfo_Request();
    return message;
  },
};

function createBasePushServiceServerInfo_Reply(): PushServiceServerInfo_Reply {
  return { publicKey: new Uint8Array(), supportedTokenTypes: [] };
}

export const PushServiceServerInfo_Reply = {
  encode(message: PushServiceServerInfo_Reply, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.publicKey.length !== 0) {
      writer.uint32(10).bytes(message.publicKey);
    }
    for (const v of message.supportedTokenTypes) {
      PushServiceSupportedTokenType.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PushServiceServerInfo_Reply {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePushServiceServerInfo_Reply();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.publicKey = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.supportedTokenTypes.push(PushServiceSupportedTokenType.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PushServiceServerInfo_Reply {
    return {
      publicKey: isSet(object.publicKey) ? bytesFromBase64(object.publicKey) : new Uint8Array(),
      supportedTokenTypes: Array.isArray(object?.supportedTokenTypes)
        ? object.supportedTokenTypes.map((e: any) => PushServiceSupportedTokenType.fromJSON(e))
        : [],
    };
  },

  toJSON(message: PushServiceServerInfo_Reply): unknown {
    const obj: any = {};
    message.publicKey !== undefined &&
      (obj.publicKey = base64FromBytes(message.publicKey !== undefined ? message.publicKey : new Uint8Array()));
    if (message.supportedTokenTypes) {
      obj.supportedTokenTypes = message.supportedTokenTypes.map((e) =>
        e ? PushServiceSupportedTokenType.toJSON(e) : undefined
      );
    } else {
      obj.supportedTokenTypes = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PushServiceServerInfo_Reply>, I>>(base?: I): PushServiceServerInfo_Reply {
    return PushServiceServerInfo_Reply.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<PushServiceServerInfo_Reply>, I>>(object: I): PushServiceServerInfo_Reply {
    const message = createBasePushServiceServerInfo_Reply();
    message.publicKey = object.publicKey ?? new Uint8Array();
    message.supportedTokenTypes =
      object.supportedTokenTypes?.map((e) => PushServiceSupportedTokenType.fromPartial(e)) || [];
    return message;
  },
};

function createBasePushServiceSupportedTokenType(): PushServiceSupportedTokenType {
  return { appBundleId: "", tokenType: 0 };
}

export const PushServiceSupportedTokenType = {
  encode(message: PushServiceSupportedTokenType, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.appBundleId !== "") {
      writer.uint32(10).string(message.appBundleId);
    }
    if (message.tokenType !== 0) {
      writer.uint32(16).int32(message.tokenType);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PushServiceSupportedTokenType {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePushServiceSupportedTokenType();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.appBundleId = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.tokenType = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PushServiceSupportedTokenType {
    return {
      appBundleId: isSet(object.appBundleId) ? String(object.appBundleId) : "",
      tokenType: isSet(object.tokenType) ? pushServiceTokenTypeFromJSON(object.tokenType) : 0,
    };
  },

  toJSON(message: PushServiceSupportedTokenType): unknown {
    const obj: any = {};
    message.appBundleId !== undefined && (obj.appBundleId = message.appBundleId);
    message.tokenType !== undefined && (obj.tokenType = pushServiceTokenTypeToJSON(message.tokenType));
    return obj;
  },

  create<I extends Exact<DeepPartial<PushServiceSupportedTokenType>, I>>(base?: I): PushServiceSupportedTokenType {
    return PushServiceSupportedTokenType.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<PushServiceSupportedTokenType>, I>>(
    object: I,
  ): PushServiceSupportedTokenType {
    const message = createBasePushServiceSupportedTokenType();
    message.appBundleId = object.appBundleId ?? "";
    message.tokenType = object.tokenType ?? 0;
    return message;
  },
};

function createBasePushServiceSend(): PushServiceSend {
  return {};
}

export const PushServiceSend = {
  encode(_: PushServiceSend, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PushServiceSend {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePushServiceSend();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): PushServiceSend {
    return {};
  },

  toJSON(_: PushServiceSend): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<PushServiceSend>, I>>(base?: I): PushServiceSend {
    return PushServiceSend.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<PushServiceSend>, I>>(_: I): PushServiceSend {
    const message = createBasePushServiceSend();
    return message;
  },
};

function createBasePushServiceSend_Request(): PushServiceSend_Request {
  return { envelope: undefined, priority: 0, receivers: [] };
}

export const PushServiceSend_Request = {
  encode(message: PushServiceSend_Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.envelope !== undefined) {
      OutOfStoreMessageEnvelope.encode(message.envelope, writer.uint32(10).fork()).ldelim();
    }
    if (message.priority !== 0) {
      writer.uint32(16).int32(message.priority);
    }
    for (const v of message.receivers) {
      PushServiceOpaqueReceiver.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PushServiceSend_Request {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePushServiceSend_Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.envelope = OutOfStoreMessageEnvelope.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.priority = reader.int32() as any;
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.receivers.push(PushServiceOpaqueReceiver.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PushServiceSend_Request {
    return {
      envelope: isSet(object.envelope) ? OutOfStoreMessageEnvelope.fromJSON(object.envelope) : undefined,
      priority: isSet(object.priority) ? pushServicePriorityFromJSON(object.priority) : 0,
      receivers: Array.isArray(object?.receivers)
        ? object.receivers.map((e: any) => PushServiceOpaqueReceiver.fromJSON(e))
        : [],
    };
  },

  toJSON(message: PushServiceSend_Request): unknown {
    const obj: any = {};
    message.envelope !== undefined &&
      (obj.envelope = message.envelope ? OutOfStoreMessageEnvelope.toJSON(message.envelope) : undefined);
    message.priority !== undefined && (obj.priority = pushServicePriorityToJSON(message.priority));
    if (message.receivers) {
      obj.receivers = message.receivers.map((e) => e ? PushServiceOpaqueReceiver.toJSON(e) : undefined);
    } else {
      obj.receivers = [];
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<PushServiceSend_Request>, I>>(base?: I): PushServiceSend_Request {
    return PushServiceSend_Request.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<PushServiceSend_Request>, I>>(object: I): PushServiceSend_Request {
    const message = createBasePushServiceSend_Request();
    message.envelope = (object.envelope !== undefined && object.envelope !== null)
      ? OutOfStoreMessageEnvelope.fromPartial(object.envelope)
      : undefined;
    message.priority = object.priority ?? 0;
    message.receivers = object.receivers?.map((e) => PushServiceOpaqueReceiver.fromPartial(e)) || [];
    return message;
  },
};

function createBasePushServiceSend_Reply(): PushServiceSend_Reply {
  return {};
}

export const PushServiceSend_Reply = {
  encode(_: PushServiceSend_Reply, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PushServiceSend_Reply {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePushServiceSend_Reply();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(_: any): PushServiceSend_Reply {
    return {};
  },

  toJSON(_: PushServiceSend_Reply): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<PushServiceSend_Reply>, I>>(base?: I): PushServiceSend_Reply {
    return PushServiceSend_Reply.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<PushServiceSend_Reply>, I>>(_: I): PushServiceSend_Reply {
    const message = createBasePushServiceSend_Reply();
    return message;
  },
};

function createBaseOutOfStoreMessageEnvelope(): OutOfStoreMessageEnvelope {
  return { nonce: new Uint8Array(), box: new Uint8Array(), groupReference: new Uint8Array() };
}

export const OutOfStoreMessageEnvelope = {
  encode(message: OutOfStoreMessageEnvelope, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nonce.length !== 0) {
      writer.uint32(10).bytes(message.nonce);
    }
    if (message.box.length !== 0) {
      writer.uint32(18).bytes(message.box);
    }
    if (message.groupReference.length !== 0) {
      writer.uint32(34).bytes(message.groupReference);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OutOfStoreMessageEnvelope {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOutOfStoreMessageEnvelope();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.nonce = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.box = reader.bytes();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.groupReference = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OutOfStoreMessageEnvelope {
    return {
      nonce: isSet(object.nonce) ? bytesFromBase64(object.nonce) : new Uint8Array(),
      box: isSet(object.box) ? bytesFromBase64(object.box) : new Uint8Array(),
      groupReference: isSet(object.groupReference) ? bytesFromBase64(object.groupReference) : new Uint8Array(),
    };
  },

  toJSON(message: OutOfStoreMessageEnvelope): unknown {
    const obj: any = {};
    message.nonce !== undefined &&
      (obj.nonce = base64FromBytes(message.nonce !== undefined ? message.nonce : new Uint8Array()));
    message.box !== undefined &&
      (obj.box = base64FromBytes(message.box !== undefined ? message.box : new Uint8Array()));
    message.groupReference !== undefined &&
      (obj.groupReference = base64FromBytes(
        message.groupReference !== undefined ? message.groupReference : new Uint8Array(),
      ));
    return obj;
  },

  create<I extends Exact<DeepPartial<OutOfStoreMessageEnvelope>, I>>(base?: I): OutOfStoreMessageEnvelope {
    return OutOfStoreMessageEnvelope.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<OutOfStoreMessageEnvelope>, I>>(object: I): OutOfStoreMessageEnvelope {
    const message = createBaseOutOfStoreMessageEnvelope();
    message.nonce = object.nonce ?? new Uint8Array();
    message.box = object.box ?? new Uint8Array();
    message.groupReference = object.groupReference ?? new Uint8Array();
    return message;
  },
};

function createBaseOutOfStoreExposedData(): OutOfStoreExposedData {
  return { nonce: new Uint8Array(), box: new Uint8Array() };
}

export const OutOfStoreExposedData = {
  encode(message: OutOfStoreExposedData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.nonce.length !== 0) {
      writer.uint32(10).bytes(message.nonce);
    }
    if (message.box.length !== 0) {
      writer.uint32(18).bytes(message.box);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): OutOfStoreExposedData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseOutOfStoreExposedData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.nonce = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.box = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): OutOfStoreExposedData {
    return {
      nonce: isSet(object.nonce) ? bytesFromBase64(object.nonce) : new Uint8Array(),
      box: isSet(object.box) ? bytesFromBase64(object.box) : new Uint8Array(),
    };
  },

  toJSON(message: OutOfStoreExposedData): unknown {
    const obj: any = {};
    message.nonce !== undefined &&
      (obj.nonce = base64FromBytes(message.nonce !== undefined ? message.nonce : new Uint8Array()));
    message.box !== undefined &&
      (obj.box = base64FromBytes(message.box !== undefined ? message.box : new Uint8Array()));
    return obj;
  },

  create<I extends Exact<DeepPartial<OutOfStoreExposedData>, I>>(base?: I): OutOfStoreExposedData {
    return OutOfStoreExposedData.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<OutOfStoreExposedData>, I>>(object: I): OutOfStoreExposedData {
    const message = createBaseOutOfStoreExposedData();
    message.nonce = object.nonce ?? new Uint8Array();
    message.box = object.box ?? new Uint8Array();
    return message;
  },
};

function createBasePushServiceOpaqueReceiver(): PushServiceOpaqueReceiver {
  return { opaqueToken: new Uint8Array(), serviceAddr: "" };
}

export const PushServiceOpaqueReceiver = {
  encode(message: PushServiceOpaqueReceiver, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.opaqueToken.length !== 0) {
      writer.uint32(10).bytes(message.opaqueToken);
    }
    if (message.serviceAddr !== "") {
      writer.uint32(18).string(message.serviceAddr);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): PushServiceOpaqueReceiver {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePushServiceOpaqueReceiver();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.opaqueToken = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.serviceAddr = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): PushServiceOpaqueReceiver {
    return {
      opaqueToken: isSet(object.opaqueToken) ? bytesFromBase64(object.opaqueToken) : new Uint8Array(),
      serviceAddr: isSet(object.serviceAddr) ? String(object.serviceAddr) : "",
    };
  },

  toJSON(message: PushServiceOpaqueReceiver): unknown {
    const obj: any = {};
    message.opaqueToken !== undefined &&
      (obj.opaqueToken = base64FromBytes(message.opaqueToken !== undefined ? message.opaqueToken : new Uint8Array()));
    message.serviceAddr !== undefined && (obj.serviceAddr = message.serviceAddr);
    return obj;
  },

  create<I extends Exact<DeepPartial<PushServiceOpaqueReceiver>, I>>(base?: I): PushServiceOpaqueReceiver {
    return PushServiceOpaqueReceiver.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<PushServiceOpaqueReceiver>, I>>(object: I): PushServiceOpaqueReceiver {
    const message = createBasePushServiceOpaqueReceiver();
    message.opaqueToken = object.opaqueToken ?? new Uint8Array();
    message.serviceAddr = object.serviceAddr ?? "";
    return message;
  },
};

function createBaseDecryptedPush(): DecryptedPush {
  return {
    accountId: "",
    accountName: "",
    conversationPublicKey: "",
    conversationDisplayName: "",
    memberPublicKey: "",
    memberDisplayName: "",
    pushType: 0,
    payloadAttrsJson: "",
    deepLink: "",
    alreadyReceived: false,
    accountMuted: false,
    conversationMuted: false,
    hidePreview: false,
  };
}

export const DecryptedPush = {
  encode(message: DecryptedPush, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accountId !== "") {
      writer.uint32(10).string(message.accountId);
    }
    if (message.accountName !== "") {
      writer.uint32(18).string(message.accountName);
    }
    if (message.conversationPublicKey !== "") {
      writer.uint32(26).string(message.conversationPublicKey);
    }
    if (message.conversationDisplayName !== "") {
      writer.uint32(34).string(message.conversationDisplayName);
    }
    if (message.memberPublicKey !== "") {
      writer.uint32(42).string(message.memberPublicKey);
    }
    if (message.memberDisplayName !== "") {
      writer.uint32(50).string(message.memberDisplayName);
    }
    if (message.pushType !== 0) {
      writer.uint32(56).int32(message.pushType);
    }
    if (message.payloadAttrsJson !== "") {
      writer.uint32(66).string(message.payloadAttrsJson);
    }
    if (message.deepLink !== "") {
      writer.uint32(74).string(message.deepLink);
    }
    if (message.alreadyReceived === true) {
      writer.uint32(80).bool(message.alreadyReceived);
    }
    if (message.accountMuted === true) {
      writer.uint32(88).bool(message.accountMuted);
    }
    if (message.conversationMuted === true) {
      writer.uint32(96).bool(message.conversationMuted);
    }
    if (message.hidePreview === true) {
      writer.uint32(104).bool(message.hidePreview);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DecryptedPush {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDecryptedPush();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.accountId = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.accountName = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.conversationPublicKey = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.conversationDisplayName = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.memberPublicKey = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.memberDisplayName = reader.string();
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.pushType = reader.int32() as any;
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.payloadAttrsJson = reader.string();
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.deepLink = reader.string();
          continue;
        case 10:
          if (tag !== 80) {
            break;
          }

          message.alreadyReceived = reader.bool();
          continue;
        case 11:
          if (tag !== 88) {
            break;
          }

          message.accountMuted = reader.bool();
          continue;
        case 12:
          if (tag !== 96) {
            break;
          }

          message.conversationMuted = reader.bool();
          continue;
        case 13:
          if (tag !== 104) {
            break;
          }

          message.hidePreview = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DecryptedPush {
    return {
      accountId: isSet(object.accountId) ? String(object.accountId) : "",
      accountName: isSet(object.accountName) ? String(object.accountName) : "",
      conversationPublicKey: isSet(object.conversationPublicKey) ? String(object.conversationPublicKey) : "",
      conversationDisplayName: isSet(object.conversationDisplayName) ? String(object.conversationDisplayName) : "",
      memberPublicKey: isSet(object.memberPublicKey) ? String(object.memberPublicKey) : "",
      memberDisplayName: isSet(object.memberDisplayName) ? String(object.memberDisplayName) : "",
      pushType: isSet(object.pushType) ? decryptedPush_PushTypeFromJSON(object.pushType) : 0,
      payloadAttrsJson: isSet(object.payloadAttrsJson) ? String(object.payloadAttrsJson) : "",
      deepLink: isSet(object.deepLink) ? String(object.deepLink) : "",
      alreadyReceived: isSet(object.alreadyReceived) ? Boolean(object.alreadyReceived) : false,
      accountMuted: isSet(object.accountMuted) ? Boolean(object.accountMuted) : false,
      conversationMuted: isSet(object.conversationMuted) ? Boolean(object.conversationMuted) : false,
      hidePreview: isSet(object.hidePreview) ? Boolean(object.hidePreview) : false,
    };
  },

  toJSON(message: DecryptedPush): unknown {
    const obj: any = {};
    message.accountId !== undefined && (obj.accountId = message.accountId);
    message.accountName !== undefined && (obj.accountName = message.accountName);
    message.conversationPublicKey !== undefined && (obj.conversationPublicKey = message.conversationPublicKey);
    message.conversationDisplayName !== undefined && (obj.conversationDisplayName = message.conversationDisplayName);
    message.memberPublicKey !== undefined && (obj.memberPublicKey = message.memberPublicKey);
    message.memberDisplayName !== undefined && (obj.memberDisplayName = message.memberDisplayName);
    message.pushType !== undefined && (obj.pushType = decryptedPush_PushTypeToJSON(message.pushType));
    message.payloadAttrsJson !== undefined && (obj.payloadAttrsJson = message.payloadAttrsJson);
    message.deepLink !== undefined && (obj.deepLink = message.deepLink);
    message.alreadyReceived !== undefined && (obj.alreadyReceived = message.alreadyReceived);
    message.accountMuted !== undefined && (obj.accountMuted = message.accountMuted);
    message.conversationMuted !== undefined && (obj.conversationMuted = message.conversationMuted);
    message.hidePreview !== undefined && (obj.hidePreview = message.hidePreview);
    return obj;
  },

  create<I extends Exact<DeepPartial<DecryptedPush>, I>>(base?: I): DecryptedPush {
    return DecryptedPush.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<DecryptedPush>, I>>(object: I): DecryptedPush {
    const message = createBaseDecryptedPush();
    message.accountId = object.accountId ?? "";
    message.accountName = object.accountName ?? "";
    message.conversationPublicKey = object.conversationPublicKey ?? "";
    message.conversationDisplayName = object.conversationDisplayName ?? "";
    message.memberPublicKey = object.memberPublicKey ?? "";
    message.memberDisplayName = object.memberDisplayName ?? "";
    message.pushType = object.pushType ?? 0;
    message.payloadAttrsJson = object.payloadAttrsJson ?? "";
    message.deepLink = object.deepLink ?? "";
    message.alreadyReceived = object.alreadyReceived ?? false;
    message.accountMuted = object.accountMuted ?? false;
    message.conversationMuted = object.conversationMuted ?? false;
    message.hidePreview = object.hidePreview ?? false;
    return message;
  },
};

function createBaseFormatedPush(): FormatedPush {
  return {
    pushType: 0,
    title: "",
    subtitle: "",
    body: "",
    deepLink: "",
    muted: false,
    hidePreview: false,
    conversationIdentifier: "",
  };
}

export const FormatedPush = {
  encode(message: FormatedPush, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pushType !== 0) {
      writer.uint32(8).int32(message.pushType);
    }
    if (message.title !== "") {
      writer.uint32(18).string(message.title);
    }
    if (message.subtitle !== "") {
      writer.uint32(26).string(message.subtitle);
    }
    if (message.body !== "") {
      writer.uint32(34).string(message.body);
    }
    if (message.deepLink !== "") {
      writer.uint32(42).string(message.deepLink);
    }
    if (message.muted === true) {
      writer.uint32(48).bool(message.muted);
    }
    if (message.hidePreview === true) {
      writer.uint32(56).bool(message.hidePreview);
    }
    if (message.conversationIdentifier !== "") {
      writer.uint32(66).string(message.conversationIdentifier);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FormatedPush {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFormatedPush();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.pushType = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.title = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.subtitle = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.body = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.deepLink = reader.string();
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.muted = reader.bool();
          continue;
        case 7:
          if (tag !== 56) {
            break;
          }

          message.hidePreview = reader.bool();
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.conversationIdentifier = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FormatedPush {
    return {
      pushType: isSet(object.pushType) ? decryptedPush_PushTypeFromJSON(object.pushType) : 0,
      title: isSet(object.title) ? String(object.title) : "",
      subtitle: isSet(object.subtitle) ? String(object.subtitle) : "",
      body: isSet(object.body) ? String(object.body) : "",
      deepLink: isSet(object.deepLink) ? String(object.deepLink) : "",
      muted: isSet(object.muted) ? Boolean(object.muted) : false,
      hidePreview: isSet(object.hidePreview) ? Boolean(object.hidePreview) : false,
      conversationIdentifier: isSet(object.conversationIdentifier) ? String(object.conversationIdentifier) : "",
    };
  },

  toJSON(message: FormatedPush): unknown {
    const obj: any = {};
    message.pushType !== undefined && (obj.pushType = decryptedPush_PushTypeToJSON(message.pushType));
    message.title !== undefined && (obj.title = message.title);
    message.subtitle !== undefined && (obj.subtitle = message.subtitle);
    message.body !== undefined && (obj.body = message.body);
    message.deepLink !== undefined && (obj.deepLink = message.deepLink);
    message.muted !== undefined && (obj.muted = message.muted);
    message.hidePreview !== undefined && (obj.hidePreview = message.hidePreview);
    message.conversationIdentifier !== undefined && (obj.conversationIdentifier = message.conversationIdentifier);
    return obj;
  },

  create<I extends Exact<DeepPartial<FormatedPush>, I>>(base?: I): FormatedPush {
    return FormatedPush.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<FormatedPush>, I>>(object: I): FormatedPush {
    const message = createBaseFormatedPush();
    message.pushType = object.pushType ?? 0;
    message.title = object.title ?? "";
    message.subtitle = object.subtitle ?? "";
    message.body = object.body ?? "";
    message.deepLink = object.deepLink ?? "";
    message.muted = object.muted ?? false;
    message.hidePreview = object.hidePreview ?? false;
    message.conversationIdentifier = object.conversationIdentifier ?? "";
    return message;
  },
};

/** PushService */
export interface PushService {
  /** ServerInfo retrieves metadata about the current push service */
  ServerInfo(
    request: DeepPartial<PushServiceServerInfo_Request>,
    metadata?: grpc.Metadata,
  ): Promise<PushServiceServerInfo_Reply>;
  /** Send dispatch a push payload to one or multiple recipients */
  Send(request: DeepPartial<PushServiceSend_Request>, metadata?: grpc.Metadata): Promise<PushServiceSend_Reply>;
}

export class PushServiceClientImpl implements PushService {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.ServerInfo = this.ServerInfo.bind(this);
    this.Send = this.Send.bind(this);
  }

  ServerInfo(
    request: DeepPartial<PushServiceServerInfo_Request>,
    metadata?: grpc.Metadata,
  ): Promise<PushServiceServerInfo_Reply> {
    return this.rpc.unary(PushServiceServerInfoDesc, PushServiceServerInfo_Request.fromPartial(request), metadata);
  }

  Send(request: DeepPartial<PushServiceSend_Request>, metadata?: grpc.Metadata): Promise<PushServiceSend_Reply> {
    return this.rpc.unary(PushServiceSendDesc, PushServiceSend_Request.fromPartial(request), metadata);
  }
}

export const PushServiceDesc = { serviceName: "weshnet.push.v1.PushService" };

export const PushServiceServerInfoDesc: UnaryMethodDefinitionish = {
  methodName: "ServerInfo",
  service: PushServiceDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return PushServiceServerInfo_Request.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      const value = PushServiceServerInfo_Reply.decode(data);
      return {
        ...value,
        toObject() {
          return value;
        },
      };
    },
  } as any,
};

export const PushServiceSendDesc: UnaryMethodDefinitionish = {
  methodName: "Send",
  service: PushServiceDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return PushServiceSend_Request.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      const value = PushServiceSend_Reply.decode(data);
      return {
        ...value,
        toObject() {
          return value;
        },
      };
    },
  } as any,
};

interface UnaryMethodDefinitionishR extends grpc.UnaryMethodDefinition<any, any> {
  requestStream: any;
  responseStream: any;
}

type UnaryMethodDefinitionish = UnaryMethodDefinitionishR;

interface Rpc {
  unary<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    request: any,
    metadata: grpc.Metadata | undefined,
  ): Promise<any>;
}

export class GrpcWebImpl {
  private host: string;
  private options: {
    transport?: grpc.TransportFactory;

    debug?: boolean;
    metadata?: grpc.Metadata;
    upStreamRetryCodes?: number[];
  };

  constructor(
    host: string,
    options: {
      transport?: grpc.TransportFactory;

      debug?: boolean;
      metadata?: grpc.Metadata;
      upStreamRetryCodes?: number[];
    },
  ) {
    this.host = host;
    this.options = options;
  }

  unary<T extends UnaryMethodDefinitionish>(
    methodDesc: T,
    _request: any,
    metadata: grpc.Metadata | undefined,
  ): Promise<any> {
    const request = { ..._request, ...methodDesc.requestType };
    const maybeCombinedMetadata = metadata && this.options.metadata
      ? new BrowserHeaders({ ...this.options?.metadata.headersMap, ...metadata?.headersMap })
      : metadata || this.options.metadata;
    return new Promise((resolve, reject) => {
      grpc.unary(methodDesc, {
        request,
        host: this.host,
        metadata: maybeCombinedMetadata,
        transport: this.options.transport,
        debug: this.options.debug,
        onEnd: function (response) {
          if (response.status === grpc.Code.OK) {
            resolve(response.message!.toObject());
          } else {
            const err = new GrpcWebError(response.statusMessage, response.status, response.trailers);
            reject(err);
          }
        },
      });
    });
  }
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

export class GrpcWebError extends tsProtoGlobalThis.Error {
  constructor(message: string, public code: grpc.Code, public metadata: grpc.Metadata) {
    super(message);
  }
}
