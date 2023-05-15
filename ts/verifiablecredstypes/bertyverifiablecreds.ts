/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "weshnet.account.v1";

export enum FlowType {
  FlowTypeUndefined = 0,
  /** FlowTypeCode - FlowTypeCode asks users a code sent on a side channel */
  FlowTypeCode = 1,
  /** FlowTypeAuth - FlowTypeAuth currently unimplemented */
  FlowTypeAuth = 2,
  /** FlowTypeProof - FlowTypeProof currently unimplemented */
  FlowTypeProof = 3,
  UNRECOGNIZED = -1,
}

export function flowTypeFromJSON(object: any): FlowType {
  switch (object) {
    case 0:
    case "FlowTypeUndefined":
      return FlowType.FlowTypeUndefined;
    case 1:
    case "FlowTypeCode":
      return FlowType.FlowTypeCode;
    case 2:
    case "FlowTypeAuth":
      return FlowType.FlowTypeAuth;
    case 3:
    case "FlowTypeProof":
      return FlowType.FlowTypeProof;
    case -1:
    case "UNRECOGNIZED":
    default:
      return FlowType.UNRECOGNIZED;
  }
}

export function flowTypeToJSON(object: FlowType): string {
  switch (object) {
    case FlowType.FlowTypeUndefined:
      return "FlowTypeUndefined";
    case FlowType.FlowTypeCode:
      return "FlowTypeCode";
    case FlowType.FlowTypeAuth:
      return "FlowTypeAuth";
    case FlowType.FlowTypeProof:
      return "FlowTypeProof";
    case FlowType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum CodeStrategy {
  CodeStrategyUndefined = 0,
  /** CodeStrategy6Digits - CodeStrategy6Digits currently unimplemented */
  CodeStrategy6Digits = 1,
  /** CodeStrategy10Chars - CodeStrategy10Chars currently unimplemented */
  CodeStrategy10Chars = 2,
  /** CodeStrategyMocked6Zeroes - CodeStrategyMocked6Zeroes must only be used in testing */
  CodeStrategyMocked6Zeroes = 999,
  UNRECOGNIZED = -1,
}

export function codeStrategyFromJSON(object: any): CodeStrategy {
  switch (object) {
    case 0:
    case "CodeStrategyUndefined":
      return CodeStrategy.CodeStrategyUndefined;
    case 1:
    case "CodeStrategy6Digits":
      return CodeStrategy.CodeStrategy6Digits;
    case 2:
    case "CodeStrategy10Chars":
      return CodeStrategy.CodeStrategy10Chars;
    case 999:
    case "CodeStrategyMocked6Zeroes":
      return CodeStrategy.CodeStrategyMocked6Zeroes;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CodeStrategy.UNRECOGNIZED;
  }
}

export function codeStrategyToJSON(object: CodeStrategy): string {
  switch (object) {
    case CodeStrategy.CodeStrategyUndefined:
      return "CodeStrategyUndefined";
    case CodeStrategy.CodeStrategy6Digits:
      return "CodeStrategy6Digits";
    case CodeStrategy.CodeStrategy10Chars:
      return "CodeStrategy10Chars";
    case CodeStrategy.CodeStrategyMocked6Zeroes:
      return "CodeStrategyMocked6Zeroes";
    case CodeStrategy.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** StateChallenge serialized and signed state used when requesting a challenge */
export interface StateChallenge {
  timestamp: Uint8Array;
  nonce: Uint8Array;
  bertyLink: string;
  redirectUri: string;
  state: string;
}

/** StateCode serialized and signed state used when requesting a code */
export interface StateCode {
  timestamp: Uint8Array;
  bertyLink: string;
  codeStrategy: CodeStrategy;
  identifier: string;
  code: string;
  redirectUri: string;
  state: string;
}

export interface AccountCryptoChallenge {
  challenge: string;
}

function createBaseStateChallenge(): StateChallenge {
  return { timestamp: new Uint8Array(), nonce: new Uint8Array(), bertyLink: "", redirectUri: "", state: "" };
}

export const StateChallenge = {
  encode(message: StateChallenge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timestamp.length !== 0) {
      writer.uint32(10).bytes(message.timestamp);
    }
    if (message.nonce.length !== 0) {
      writer.uint32(18).bytes(message.nonce);
    }
    if (message.bertyLink !== "") {
      writer.uint32(26).string(message.bertyLink);
    }
    if (message.redirectUri !== "") {
      writer.uint32(34).string(message.redirectUri);
    }
    if (message.state !== "") {
      writer.uint32(42).string(message.state);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StateChallenge {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStateChallenge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.timestamp = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.nonce = reader.bytes();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.bertyLink = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.redirectUri = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.state = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StateChallenge {
    return {
      timestamp: isSet(object.timestamp) ? bytesFromBase64(object.timestamp) : new Uint8Array(),
      nonce: isSet(object.nonce) ? bytesFromBase64(object.nonce) : new Uint8Array(),
      bertyLink: isSet(object.bertyLink) ? String(object.bertyLink) : "",
      redirectUri: isSet(object.redirectUri) ? String(object.redirectUri) : "",
      state: isSet(object.state) ? String(object.state) : "",
    };
  },

  toJSON(message: StateChallenge): unknown {
    const obj: any = {};
    message.timestamp !== undefined &&
      (obj.timestamp = base64FromBytes(message.timestamp !== undefined ? message.timestamp : new Uint8Array()));
    message.nonce !== undefined &&
      (obj.nonce = base64FromBytes(message.nonce !== undefined ? message.nonce : new Uint8Array()));
    message.bertyLink !== undefined && (obj.bertyLink = message.bertyLink);
    message.redirectUri !== undefined && (obj.redirectUri = message.redirectUri);
    message.state !== undefined && (obj.state = message.state);
    return obj;
  },

  create<I extends Exact<DeepPartial<StateChallenge>, I>>(base?: I): StateChallenge {
    return StateChallenge.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<StateChallenge>, I>>(object: I): StateChallenge {
    const message = createBaseStateChallenge();
    message.timestamp = object.timestamp ?? new Uint8Array();
    message.nonce = object.nonce ?? new Uint8Array();
    message.bertyLink = object.bertyLink ?? "";
    message.redirectUri = object.redirectUri ?? "";
    message.state = object.state ?? "";
    return message;
  },
};

function createBaseStateCode(): StateCode {
  return {
    timestamp: new Uint8Array(),
    bertyLink: "",
    codeStrategy: 0,
    identifier: "",
    code: "",
    redirectUri: "",
    state: "",
  };
}

export const StateCode = {
  encode(message: StateCode, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.timestamp.length !== 0) {
      writer.uint32(10).bytes(message.timestamp);
    }
    if (message.bertyLink !== "") {
      writer.uint32(18).string(message.bertyLink);
    }
    if (message.codeStrategy !== 0) {
      writer.uint32(24).int32(message.codeStrategy);
    }
    if (message.identifier !== "") {
      writer.uint32(34).string(message.identifier);
    }
    if (message.code !== "") {
      writer.uint32(42).string(message.code);
    }
    if (message.redirectUri !== "") {
      writer.uint32(50).string(message.redirectUri);
    }
    if (message.state !== "") {
      writer.uint32(58).string(message.state);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StateCode {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStateCode();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.timestamp = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.bertyLink = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.codeStrategy = reader.int32() as any;
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.identifier = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.code = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.redirectUri = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.state = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): StateCode {
    return {
      timestamp: isSet(object.timestamp) ? bytesFromBase64(object.timestamp) : new Uint8Array(),
      bertyLink: isSet(object.bertyLink) ? String(object.bertyLink) : "",
      codeStrategy: isSet(object.codeStrategy) ? codeStrategyFromJSON(object.codeStrategy) : 0,
      identifier: isSet(object.identifier) ? String(object.identifier) : "",
      code: isSet(object.code) ? String(object.code) : "",
      redirectUri: isSet(object.redirectUri) ? String(object.redirectUri) : "",
      state: isSet(object.state) ? String(object.state) : "",
    };
  },

  toJSON(message: StateCode): unknown {
    const obj: any = {};
    message.timestamp !== undefined &&
      (obj.timestamp = base64FromBytes(message.timestamp !== undefined ? message.timestamp : new Uint8Array()));
    message.bertyLink !== undefined && (obj.bertyLink = message.bertyLink);
    message.codeStrategy !== undefined && (obj.codeStrategy = codeStrategyToJSON(message.codeStrategy));
    message.identifier !== undefined && (obj.identifier = message.identifier);
    message.code !== undefined && (obj.code = message.code);
    message.redirectUri !== undefined && (obj.redirectUri = message.redirectUri);
    message.state !== undefined && (obj.state = message.state);
    return obj;
  },

  create<I extends Exact<DeepPartial<StateCode>, I>>(base?: I): StateCode {
    return StateCode.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<StateCode>, I>>(object: I): StateCode {
    const message = createBaseStateCode();
    message.timestamp = object.timestamp ?? new Uint8Array();
    message.bertyLink = object.bertyLink ?? "";
    message.codeStrategy = object.codeStrategy ?? 0;
    message.identifier = object.identifier ?? "";
    message.code = object.code ?? "";
    message.redirectUri = object.redirectUri ?? "";
    message.state = object.state ?? "";
    return message;
  },
};

function createBaseAccountCryptoChallenge(): AccountCryptoChallenge {
  return { challenge: "" };
}

export const AccountCryptoChallenge = {
  encode(message: AccountCryptoChallenge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.challenge !== "") {
      writer.uint32(10).string(message.challenge);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AccountCryptoChallenge {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAccountCryptoChallenge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.challenge = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AccountCryptoChallenge {
    return { challenge: isSet(object.challenge) ? String(object.challenge) : "" };
  },

  toJSON(message: AccountCryptoChallenge): unknown {
    const obj: any = {};
    message.challenge !== undefined && (obj.challenge = message.challenge);
    return obj;
  },

  create<I extends Exact<DeepPartial<AccountCryptoChallenge>, I>>(base?: I): AccountCryptoChallenge {
    return AccountCryptoChallenge.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<AccountCryptoChallenge>, I>>(object: I): AccountCryptoChallenge {
    const message = createBaseAccountCryptoChallenge();
    message.challenge = object.challenge ?? "";
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
