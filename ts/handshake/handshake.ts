/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "handshake";

export interface BoxEnvelope {
  box: Uint8Array;
}

export interface HelloPayload {
  ephemeralPubKey: Uint8Array;
}

export interface RequesterAuthenticatePayload {
  requesterAccountId: Uint8Array;
  requesterAccountSig: Uint8Array;
}

export interface ResponderAcceptPayload {
  responderAccountSig: Uint8Array;
}

export interface RequesterAcknowledgePayload {
  success: boolean;
}

function createBaseBoxEnvelope(): BoxEnvelope {
  return { box: new Uint8Array() };
}

export const BoxEnvelope = {
  encode(message: BoxEnvelope, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.box.length !== 0) {
      writer.uint32(10).bytes(message.box);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BoxEnvelope {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBoxEnvelope();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
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

  fromJSON(object: any): BoxEnvelope {
    return { box: isSet(object.box) ? bytesFromBase64(object.box) : new Uint8Array() };
  },

  toJSON(message: BoxEnvelope): unknown {
    const obj: any = {};
    message.box !== undefined &&
      (obj.box = base64FromBytes(message.box !== undefined ? message.box : new Uint8Array()));
    return obj;
  },

  create<I extends Exact<DeepPartial<BoxEnvelope>, I>>(base?: I): BoxEnvelope {
    return BoxEnvelope.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<BoxEnvelope>, I>>(object: I): BoxEnvelope {
    const message = createBaseBoxEnvelope();
    message.box = object.box ?? new Uint8Array();
    return message;
  },
};

function createBaseHelloPayload(): HelloPayload {
  return { ephemeralPubKey: new Uint8Array() };
}

export const HelloPayload = {
  encode(message: HelloPayload, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ephemeralPubKey.length !== 0) {
      writer.uint32(10).bytes(message.ephemeralPubKey);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HelloPayload {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHelloPayload();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.ephemeralPubKey = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HelloPayload {
    return {
      ephemeralPubKey: isSet(object.ephemeralPubKey) ? bytesFromBase64(object.ephemeralPubKey) : new Uint8Array(),
    };
  },

  toJSON(message: HelloPayload): unknown {
    const obj: any = {};
    message.ephemeralPubKey !== undefined &&
      (obj.ephemeralPubKey = base64FromBytes(
        message.ephemeralPubKey !== undefined ? message.ephemeralPubKey : new Uint8Array(),
      ));
    return obj;
  },

  create<I extends Exact<DeepPartial<HelloPayload>, I>>(base?: I): HelloPayload {
    return HelloPayload.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<HelloPayload>, I>>(object: I): HelloPayload {
    const message = createBaseHelloPayload();
    message.ephemeralPubKey = object.ephemeralPubKey ?? new Uint8Array();
    return message;
  },
};

function createBaseRequesterAuthenticatePayload(): RequesterAuthenticatePayload {
  return { requesterAccountId: new Uint8Array(), requesterAccountSig: new Uint8Array() };
}

export const RequesterAuthenticatePayload = {
  encode(message: RequesterAuthenticatePayload, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.requesterAccountId.length !== 0) {
      writer.uint32(10).bytes(message.requesterAccountId);
    }
    if (message.requesterAccountSig.length !== 0) {
      writer.uint32(18).bytes(message.requesterAccountSig);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequesterAuthenticatePayload {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequesterAuthenticatePayload();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.requesterAccountId = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.requesterAccountSig = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequesterAuthenticatePayload {
    return {
      requesterAccountId: isSet(object.requesterAccountId)
        ? bytesFromBase64(object.requesterAccountId)
        : new Uint8Array(),
      requesterAccountSig: isSet(object.requesterAccountSig)
        ? bytesFromBase64(object.requesterAccountSig)
        : new Uint8Array(),
    };
  },

  toJSON(message: RequesterAuthenticatePayload): unknown {
    const obj: any = {};
    message.requesterAccountId !== undefined &&
      (obj.requesterAccountId = base64FromBytes(
        message.requesterAccountId !== undefined ? message.requesterAccountId : new Uint8Array(),
      ));
    message.requesterAccountSig !== undefined &&
      (obj.requesterAccountSig = base64FromBytes(
        message.requesterAccountSig !== undefined ? message.requesterAccountSig : new Uint8Array(),
      ));
    return obj;
  },

  create<I extends Exact<DeepPartial<RequesterAuthenticatePayload>, I>>(base?: I): RequesterAuthenticatePayload {
    return RequesterAuthenticatePayload.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<RequesterAuthenticatePayload>, I>>(object: I): RequesterAuthenticatePayload {
    const message = createBaseRequesterAuthenticatePayload();
    message.requesterAccountId = object.requesterAccountId ?? new Uint8Array();
    message.requesterAccountSig = object.requesterAccountSig ?? new Uint8Array();
    return message;
  },
};

function createBaseResponderAcceptPayload(): ResponderAcceptPayload {
  return { responderAccountSig: new Uint8Array() };
}

export const ResponderAcceptPayload = {
  encode(message: ResponderAcceptPayload, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.responderAccountSig.length !== 0) {
      writer.uint32(10).bytes(message.responderAccountSig);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ResponderAcceptPayload {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseResponderAcceptPayload();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.responderAccountSig = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ResponderAcceptPayload {
    return {
      responderAccountSig: isSet(object.responderAccountSig)
        ? bytesFromBase64(object.responderAccountSig)
        : new Uint8Array(),
    };
  },

  toJSON(message: ResponderAcceptPayload): unknown {
    const obj: any = {};
    message.responderAccountSig !== undefined &&
      (obj.responderAccountSig = base64FromBytes(
        message.responderAccountSig !== undefined ? message.responderAccountSig : new Uint8Array(),
      ));
    return obj;
  },

  create<I extends Exact<DeepPartial<ResponderAcceptPayload>, I>>(base?: I): ResponderAcceptPayload {
    return ResponderAcceptPayload.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ResponderAcceptPayload>, I>>(object: I): ResponderAcceptPayload {
    const message = createBaseResponderAcceptPayload();
    message.responderAccountSig = object.responderAccountSig ?? new Uint8Array();
    return message;
  },
};

function createBaseRequesterAcknowledgePayload(): RequesterAcknowledgePayload {
  return { success: false };
}

export const RequesterAcknowledgePayload = {
  encode(message: RequesterAcknowledgePayload, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.success === true) {
      writer.uint32(8).bool(message.success);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RequesterAcknowledgePayload {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRequesterAcknowledgePayload();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.success = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RequesterAcknowledgePayload {
    return { success: isSet(object.success) ? Boolean(object.success) : false };
  },

  toJSON(message: RequesterAcknowledgePayload): unknown {
    const obj: any = {};
    message.success !== undefined && (obj.success = message.success);
    return obj;
  },

  create<I extends Exact<DeepPartial<RequesterAcknowledgePayload>, I>>(base?: I): RequesterAcknowledgePayload {
    return RequesterAcknowledgePayload.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<RequesterAcknowledgePayload>, I>>(object: I): RequesterAcknowledgePayload {
    const message = createBaseRequesterAcknowledgePayload();
    message.success = object.success ?? false;
    return message;
  },
};

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
