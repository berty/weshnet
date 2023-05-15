/* eslint-disable */
import { grpc } from "@improbable-eng/grpc-web";
import { BrowserHeaders } from "browser-headers";
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Group } from "../protocoltypes";

export const protobufPackage = "weshnet.replication.v1";

export interface ReplicatedGroup {
  publicKey: string;
  signPub: string;
  linkKey: string;
  createdAt: number;
  updatedAt: number;
  metadataEntriesCount: number;
  metadataLatestHead: string;
  messageEntriesCount: number;
  messageLatestHead: string;
}

export interface ReplicatedGroupToken {
  replicatedGroupPublicKey: string;
  replicatedGroup: ReplicatedGroup | undefined;
  tokenIssuer: string;
  tokenId: string;
  createdAt: number;
}

export interface ReplicationServiceReplicateGroup {
}

export interface ReplicationServiceReplicateGroup_Request {
  group: Group | undefined;
}

export interface ReplicationServiceReplicateGroup_Reply {
  ok: boolean;
}

export interface ReplicateGlobalStats {
}

export interface ReplicateGlobalStats_Request {
}

export interface ReplicateGlobalStats_Reply {
  startedAt: number;
  replicatedGroups: number;
  totalMetadataEntries: number;
  totalMessageEntries: number;
}

export interface ReplicateGroupStats {
}

export interface ReplicateGroupStats_Request {
  groupPublicKey: string;
}

export interface ReplicateGroupStats_Reply {
  group: ReplicatedGroup | undefined;
}

function createBaseReplicatedGroup(): ReplicatedGroup {
  return {
    publicKey: "",
    signPub: "",
    linkKey: "",
    createdAt: 0,
    updatedAt: 0,
    metadataEntriesCount: 0,
    metadataLatestHead: "",
    messageEntriesCount: 0,
    messageLatestHead: "",
  };
}

export const ReplicatedGroup = {
  encode(message: ReplicatedGroup, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.publicKey !== "") {
      writer.uint32(10).string(message.publicKey);
    }
    if (message.signPub !== "") {
      writer.uint32(18).string(message.signPub);
    }
    if (message.linkKey !== "") {
      writer.uint32(26).string(message.linkKey);
    }
    if (message.createdAt !== 0) {
      writer.uint32(800).int64(message.createdAt);
    }
    if (message.updatedAt !== 0) {
      writer.uint32(808).int64(message.updatedAt);
    }
    if (message.metadataEntriesCount !== 0) {
      writer.uint32(816).int64(message.metadataEntriesCount);
    }
    if (message.metadataLatestHead !== "") {
      writer.uint32(826).string(message.metadataLatestHead);
    }
    if (message.messageEntriesCount !== 0) {
      writer.uint32(832).int64(message.messageEntriesCount);
    }
    if (message.messageLatestHead !== "") {
      writer.uint32(842).string(message.messageLatestHead);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicatedGroup {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicatedGroup();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.publicKey = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.signPub = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.linkKey = reader.string();
          continue;
        case 100:
          if (tag !== 800) {
            break;
          }

          message.createdAt = longToNumber(reader.int64() as Long);
          continue;
        case 101:
          if (tag !== 808) {
            break;
          }

          message.updatedAt = longToNumber(reader.int64() as Long);
          continue;
        case 102:
          if (tag !== 816) {
            break;
          }

          message.metadataEntriesCount = longToNumber(reader.int64() as Long);
          continue;
        case 103:
          if (tag !== 826) {
            break;
          }

          message.metadataLatestHead = reader.string();
          continue;
        case 104:
          if (tag !== 832) {
            break;
          }

          message.messageEntriesCount = longToNumber(reader.int64() as Long);
          continue;
        case 105:
          if (tag !== 842) {
            break;
          }

          message.messageLatestHead = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplicatedGroup {
    return {
      publicKey: isSet(object.publicKey) ? String(object.publicKey) : "",
      signPub: isSet(object.signPub) ? String(object.signPub) : "",
      linkKey: isSet(object.linkKey) ? String(object.linkKey) : "",
      createdAt: isSet(object.createdAt) ? Number(object.createdAt) : 0,
      updatedAt: isSet(object.updatedAt) ? Number(object.updatedAt) : 0,
      metadataEntriesCount: isSet(object.metadataEntriesCount) ? Number(object.metadataEntriesCount) : 0,
      metadataLatestHead: isSet(object.metadataLatestHead) ? String(object.metadataLatestHead) : "",
      messageEntriesCount: isSet(object.messageEntriesCount) ? Number(object.messageEntriesCount) : 0,
      messageLatestHead: isSet(object.messageLatestHead) ? String(object.messageLatestHead) : "",
    };
  },

  toJSON(message: ReplicatedGroup): unknown {
    const obj: any = {};
    message.publicKey !== undefined && (obj.publicKey = message.publicKey);
    message.signPub !== undefined && (obj.signPub = message.signPub);
    message.linkKey !== undefined && (obj.linkKey = message.linkKey);
    message.createdAt !== undefined && (obj.createdAt = Math.round(message.createdAt));
    message.updatedAt !== undefined && (obj.updatedAt = Math.round(message.updatedAt));
    message.metadataEntriesCount !== undefined && (obj.metadataEntriesCount = Math.round(message.metadataEntriesCount));
    message.metadataLatestHead !== undefined && (obj.metadataLatestHead = message.metadataLatestHead);
    message.messageEntriesCount !== undefined && (obj.messageEntriesCount = Math.round(message.messageEntriesCount));
    message.messageLatestHead !== undefined && (obj.messageLatestHead = message.messageLatestHead);
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicatedGroup>, I>>(base?: I): ReplicatedGroup {
    return ReplicatedGroup.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicatedGroup>, I>>(object: I): ReplicatedGroup {
    const message = createBaseReplicatedGroup();
    message.publicKey = object.publicKey ?? "";
    message.signPub = object.signPub ?? "";
    message.linkKey = object.linkKey ?? "";
    message.createdAt = object.createdAt ?? 0;
    message.updatedAt = object.updatedAt ?? 0;
    message.metadataEntriesCount = object.metadataEntriesCount ?? 0;
    message.metadataLatestHead = object.metadataLatestHead ?? "";
    message.messageEntriesCount = object.messageEntriesCount ?? 0;
    message.messageLatestHead = object.messageLatestHead ?? "";
    return message;
  },
};

function createBaseReplicatedGroupToken(): ReplicatedGroupToken {
  return { replicatedGroupPublicKey: "", replicatedGroup: undefined, tokenIssuer: "", tokenId: "", createdAt: 0 };
}

export const ReplicatedGroupToken = {
  encode(message: ReplicatedGroupToken, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.replicatedGroupPublicKey !== "") {
      writer.uint32(10).string(message.replicatedGroupPublicKey);
    }
    if (message.replicatedGroup !== undefined) {
      ReplicatedGroup.encode(message.replicatedGroup, writer.uint32(18).fork()).ldelim();
    }
    if (message.tokenIssuer !== "") {
      writer.uint32(26).string(message.tokenIssuer);
    }
    if (message.tokenId !== "") {
      writer.uint32(34).string(message.tokenId);
    }
    if (message.createdAt !== 0) {
      writer.uint32(40).int64(message.createdAt);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicatedGroupToken {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicatedGroupToken();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.replicatedGroupPublicKey = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.replicatedGroup = ReplicatedGroup.decode(reader, reader.uint32());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.tokenIssuer = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.tokenId = reader.string();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.createdAt = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplicatedGroupToken {
    return {
      replicatedGroupPublicKey: isSet(object.replicatedGroupPublicKey) ? String(object.replicatedGroupPublicKey) : "",
      replicatedGroup: isSet(object.replicatedGroup) ? ReplicatedGroup.fromJSON(object.replicatedGroup) : undefined,
      tokenIssuer: isSet(object.tokenIssuer) ? String(object.tokenIssuer) : "",
      tokenId: isSet(object.tokenId) ? String(object.tokenId) : "",
      createdAt: isSet(object.createdAt) ? Number(object.createdAt) : 0,
    };
  },

  toJSON(message: ReplicatedGroupToken): unknown {
    const obj: any = {};
    message.replicatedGroupPublicKey !== undefined && (obj.replicatedGroupPublicKey = message.replicatedGroupPublicKey);
    message.replicatedGroup !== undefined &&
      (obj.replicatedGroup = message.replicatedGroup ? ReplicatedGroup.toJSON(message.replicatedGroup) : undefined);
    message.tokenIssuer !== undefined && (obj.tokenIssuer = message.tokenIssuer);
    message.tokenId !== undefined && (obj.tokenId = message.tokenId);
    message.createdAt !== undefined && (obj.createdAt = Math.round(message.createdAt));
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicatedGroupToken>, I>>(base?: I): ReplicatedGroupToken {
    return ReplicatedGroupToken.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicatedGroupToken>, I>>(object: I): ReplicatedGroupToken {
    const message = createBaseReplicatedGroupToken();
    message.replicatedGroupPublicKey = object.replicatedGroupPublicKey ?? "";
    message.replicatedGroup = (object.replicatedGroup !== undefined && object.replicatedGroup !== null)
      ? ReplicatedGroup.fromPartial(object.replicatedGroup)
      : undefined;
    message.tokenIssuer = object.tokenIssuer ?? "";
    message.tokenId = object.tokenId ?? "";
    message.createdAt = object.createdAt ?? 0;
    return message;
  },
};

function createBaseReplicationServiceReplicateGroup(): ReplicationServiceReplicateGroup {
  return {};
}

export const ReplicationServiceReplicateGroup = {
  encode(_: ReplicationServiceReplicateGroup, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicationServiceReplicateGroup {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicationServiceReplicateGroup();
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

  fromJSON(_: any): ReplicationServiceReplicateGroup {
    return {};
  },

  toJSON(_: ReplicationServiceReplicateGroup): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicationServiceReplicateGroup>, I>>(
    base?: I,
  ): ReplicationServiceReplicateGroup {
    return ReplicationServiceReplicateGroup.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicationServiceReplicateGroup>, I>>(
    _: I,
  ): ReplicationServiceReplicateGroup {
    const message = createBaseReplicationServiceReplicateGroup();
    return message;
  },
};

function createBaseReplicationServiceReplicateGroup_Request(): ReplicationServiceReplicateGroup_Request {
  return { group: undefined };
}

export const ReplicationServiceReplicateGroup_Request = {
  encode(message: ReplicationServiceReplicateGroup_Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.group !== undefined) {
      Group.encode(message.group, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicationServiceReplicateGroup_Request {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicationServiceReplicateGroup_Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.group = Group.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplicationServiceReplicateGroup_Request {
    return { group: isSet(object.group) ? Group.fromJSON(object.group) : undefined };
  },

  toJSON(message: ReplicationServiceReplicateGroup_Request): unknown {
    const obj: any = {};
    message.group !== undefined && (obj.group = message.group ? Group.toJSON(message.group) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicationServiceReplicateGroup_Request>, I>>(
    base?: I,
  ): ReplicationServiceReplicateGroup_Request {
    return ReplicationServiceReplicateGroup_Request.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicationServiceReplicateGroup_Request>, I>>(
    object: I,
  ): ReplicationServiceReplicateGroup_Request {
    const message = createBaseReplicationServiceReplicateGroup_Request();
    message.group = (object.group !== undefined && object.group !== null) ? Group.fromPartial(object.group) : undefined;
    return message;
  },
};

function createBaseReplicationServiceReplicateGroup_Reply(): ReplicationServiceReplicateGroup_Reply {
  return { ok: false };
}

export const ReplicationServiceReplicateGroup_Reply = {
  encode(message: ReplicationServiceReplicateGroup_Reply, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.ok === true) {
      writer.uint32(8).bool(message.ok);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicationServiceReplicateGroup_Reply {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicationServiceReplicateGroup_Reply();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.ok = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplicationServiceReplicateGroup_Reply {
    return { ok: isSet(object.ok) ? Boolean(object.ok) : false };
  },

  toJSON(message: ReplicationServiceReplicateGroup_Reply): unknown {
    const obj: any = {};
    message.ok !== undefined && (obj.ok = message.ok);
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicationServiceReplicateGroup_Reply>, I>>(
    base?: I,
  ): ReplicationServiceReplicateGroup_Reply {
    return ReplicationServiceReplicateGroup_Reply.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicationServiceReplicateGroup_Reply>, I>>(
    object: I,
  ): ReplicationServiceReplicateGroup_Reply {
    const message = createBaseReplicationServiceReplicateGroup_Reply();
    message.ok = object.ok ?? false;
    return message;
  },
};

function createBaseReplicateGlobalStats(): ReplicateGlobalStats {
  return {};
}

export const ReplicateGlobalStats = {
  encode(_: ReplicateGlobalStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicateGlobalStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicateGlobalStats();
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

  fromJSON(_: any): ReplicateGlobalStats {
    return {};
  },

  toJSON(_: ReplicateGlobalStats): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicateGlobalStats>, I>>(base?: I): ReplicateGlobalStats {
    return ReplicateGlobalStats.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicateGlobalStats>, I>>(_: I): ReplicateGlobalStats {
    const message = createBaseReplicateGlobalStats();
    return message;
  },
};

function createBaseReplicateGlobalStats_Request(): ReplicateGlobalStats_Request {
  return {};
}

export const ReplicateGlobalStats_Request = {
  encode(_: ReplicateGlobalStats_Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicateGlobalStats_Request {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicateGlobalStats_Request();
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

  fromJSON(_: any): ReplicateGlobalStats_Request {
    return {};
  },

  toJSON(_: ReplicateGlobalStats_Request): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicateGlobalStats_Request>, I>>(base?: I): ReplicateGlobalStats_Request {
    return ReplicateGlobalStats_Request.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicateGlobalStats_Request>, I>>(_: I): ReplicateGlobalStats_Request {
    const message = createBaseReplicateGlobalStats_Request();
    return message;
  },
};

function createBaseReplicateGlobalStats_Reply(): ReplicateGlobalStats_Reply {
  return { startedAt: 0, replicatedGroups: 0, totalMetadataEntries: 0, totalMessageEntries: 0 };
}

export const ReplicateGlobalStats_Reply = {
  encode(message: ReplicateGlobalStats_Reply, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.startedAt !== 0) {
      writer.uint32(8).int64(message.startedAt);
    }
    if (message.replicatedGroups !== 0) {
      writer.uint32(16).int64(message.replicatedGroups);
    }
    if (message.totalMetadataEntries !== 0) {
      writer.uint32(24).int64(message.totalMetadataEntries);
    }
    if (message.totalMessageEntries !== 0) {
      writer.uint32(32).int64(message.totalMessageEntries);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicateGlobalStats_Reply {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicateGlobalStats_Reply();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.startedAt = longToNumber(reader.int64() as Long);
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.replicatedGroups = longToNumber(reader.int64() as Long);
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.totalMetadataEntries = longToNumber(reader.int64() as Long);
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.totalMessageEntries = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplicateGlobalStats_Reply {
    return {
      startedAt: isSet(object.startedAt) ? Number(object.startedAt) : 0,
      replicatedGroups: isSet(object.replicatedGroups) ? Number(object.replicatedGroups) : 0,
      totalMetadataEntries: isSet(object.totalMetadataEntries) ? Number(object.totalMetadataEntries) : 0,
      totalMessageEntries: isSet(object.totalMessageEntries) ? Number(object.totalMessageEntries) : 0,
    };
  },

  toJSON(message: ReplicateGlobalStats_Reply): unknown {
    const obj: any = {};
    message.startedAt !== undefined && (obj.startedAt = Math.round(message.startedAt));
    message.replicatedGroups !== undefined && (obj.replicatedGroups = Math.round(message.replicatedGroups));
    message.totalMetadataEntries !== undefined && (obj.totalMetadataEntries = Math.round(message.totalMetadataEntries));
    message.totalMessageEntries !== undefined && (obj.totalMessageEntries = Math.round(message.totalMessageEntries));
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicateGlobalStats_Reply>, I>>(base?: I): ReplicateGlobalStats_Reply {
    return ReplicateGlobalStats_Reply.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicateGlobalStats_Reply>, I>>(object: I): ReplicateGlobalStats_Reply {
    const message = createBaseReplicateGlobalStats_Reply();
    message.startedAt = object.startedAt ?? 0;
    message.replicatedGroups = object.replicatedGroups ?? 0;
    message.totalMetadataEntries = object.totalMetadataEntries ?? 0;
    message.totalMessageEntries = object.totalMessageEntries ?? 0;
    return message;
  },
};

function createBaseReplicateGroupStats(): ReplicateGroupStats {
  return {};
}

export const ReplicateGroupStats = {
  encode(_: ReplicateGroupStats, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicateGroupStats {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicateGroupStats();
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

  fromJSON(_: any): ReplicateGroupStats {
    return {};
  },

  toJSON(_: ReplicateGroupStats): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicateGroupStats>, I>>(base?: I): ReplicateGroupStats {
    return ReplicateGroupStats.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicateGroupStats>, I>>(_: I): ReplicateGroupStats {
    const message = createBaseReplicateGroupStats();
    return message;
  },
};

function createBaseReplicateGroupStats_Request(): ReplicateGroupStats_Request {
  return { groupPublicKey: "" };
}

export const ReplicateGroupStats_Request = {
  encode(message: ReplicateGroupStats_Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.groupPublicKey !== "") {
      writer.uint32(10).string(message.groupPublicKey);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicateGroupStats_Request {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicateGroupStats_Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.groupPublicKey = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplicateGroupStats_Request {
    return { groupPublicKey: isSet(object.groupPublicKey) ? String(object.groupPublicKey) : "" };
  },

  toJSON(message: ReplicateGroupStats_Request): unknown {
    const obj: any = {};
    message.groupPublicKey !== undefined && (obj.groupPublicKey = message.groupPublicKey);
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicateGroupStats_Request>, I>>(base?: I): ReplicateGroupStats_Request {
    return ReplicateGroupStats_Request.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicateGroupStats_Request>, I>>(object: I): ReplicateGroupStats_Request {
    const message = createBaseReplicateGroupStats_Request();
    message.groupPublicKey = object.groupPublicKey ?? "";
    return message;
  },
};

function createBaseReplicateGroupStats_Reply(): ReplicateGroupStats_Reply {
  return { group: undefined };
}

export const ReplicateGroupStats_Reply = {
  encode(message: ReplicateGroupStats_Reply, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.group !== undefined) {
      ReplicatedGroup.encode(message.group, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReplicateGroupStats_Reply {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReplicateGroupStats_Reply();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.group = ReplicatedGroup.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ReplicateGroupStats_Reply {
    return { group: isSet(object.group) ? ReplicatedGroup.fromJSON(object.group) : undefined };
  },

  toJSON(message: ReplicateGroupStats_Reply): unknown {
    const obj: any = {};
    message.group !== undefined && (obj.group = message.group ? ReplicatedGroup.toJSON(message.group) : undefined);
    return obj;
  },

  create<I extends Exact<DeepPartial<ReplicateGroupStats_Reply>, I>>(base?: I): ReplicateGroupStats_Reply {
    return ReplicateGroupStats_Reply.fromPartial(base ?? {});
  },

  fromPartial<I extends Exact<DeepPartial<ReplicateGroupStats_Reply>, I>>(object: I): ReplicateGroupStats_Reply {
    const message = createBaseReplicateGroupStats_Reply();
    message.group = (object.group !== undefined && object.group !== null)
      ? ReplicatedGroup.fromPartial(object.group)
      : undefined;
    return message;
  },
};

/** ReplicationService */
export interface ReplicationService {
  /** ReplicateGroup */
  ReplicateGroup(
    request: DeepPartial<ReplicationServiceReplicateGroup_Request>,
    metadata?: grpc.Metadata,
  ): Promise<ReplicationServiceReplicateGroup_Reply>;
  ReplicateGlobalStats(
    request: DeepPartial<ReplicateGlobalStats_Request>,
    metadata?: grpc.Metadata,
  ): Promise<ReplicateGlobalStats_Reply>;
  ReplicateGroupStats(
    request: DeepPartial<ReplicateGroupStats_Request>,
    metadata?: grpc.Metadata,
  ): Promise<ReplicateGroupStats_Reply>;
}

export class ReplicationServiceClientImpl implements ReplicationService {
  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.ReplicateGroup = this.ReplicateGroup.bind(this);
    this.ReplicateGlobalStats = this.ReplicateGlobalStats.bind(this);
    this.ReplicateGroupStats = this.ReplicateGroupStats.bind(this);
  }

  ReplicateGroup(
    request: DeepPartial<ReplicationServiceReplicateGroup_Request>,
    metadata?: grpc.Metadata,
  ): Promise<ReplicationServiceReplicateGroup_Reply> {
    return this.rpc.unary(
      ReplicationServiceReplicateGroupDesc,
      ReplicationServiceReplicateGroup_Request.fromPartial(request),
      metadata,
    );
  }

  ReplicateGlobalStats(
    request: DeepPartial<ReplicateGlobalStats_Request>,
    metadata?: grpc.Metadata,
  ): Promise<ReplicateGlobalStats_Reply> {
    return this.rpc.unary(
      ReplicationServiceReplicateGlobalStatsDesc,
      ReplicateGlobalStats_Request.fromPartial(request),
      metadata,
    );
  }

  ReplicateGroupStats(
    request: DeepPartial<ReplicateGroupStats_Request>,
    metadata?: grpc.Metadata,
  ): Promise<ReplicateGroupStats_Reply> {
    return this.rpc.unary(
      ReplicationServiceReplicateGroupStatsDesc,
      ReplicateGroupStats_Request.fromPartial(request),
      metadata,
    );
  }
}

export const ReplicationServiceDesc = { serviceName: "weshnet.replication.v1.ReplicationService" };

export const ReplicationServiceReplicateGroupDesc: UnaryMethodDefinitionish = {
  methodName: "ReplicateGroup",
  service: ReplicationServiceDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return ReplicationServiceReplicateGroup_Request.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      const value = ReplicationServiceReplicateGroup_Reply.decode(data);
      return {
        ...value,
        toObject() {
          return value;
        },
      };
    },
  } as any,
};

export const ReplicationServiceReplicateGlobalStatsDesc: UnaryMethodDefinitionish = {
  methodName: "ReplicateGlobalStats",
  service: ReplicationServiceDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return ReplicateGlobalStats_Request.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      const value = ReplicateGlobalStats_Reply.decode(data);
      return {
        ...value,
        toObject() {
          return value;
        },
      };
    },
  } as any,
};

export const ReplicationServiceReplicateGroupStatsDesc: UnaryMethodDefinitionish = {
  methodName: "ReplicateGroupStats",
  service: ReplicationServiceDesc,
  requestStream: false,
  responseStream: false,
  requestType: {
    serializeBinary() {
      return ReplicateGroupStats_Request.encode(this).finish();
    },
  } as any,
  responseType: {
    deserializeBinary(data: Uint8Array) {
      const value = ReplicateGroupStats_Reply.decode(data);
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

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}

export class GrpcWebError extends tsProtoGlobalThis.Error {
  constructor(message: string, public code: grpc.Code, public metadata: grpc.Metadata) {
    super(message);
  }
}
