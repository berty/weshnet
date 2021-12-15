//
//  ConnectedPeer.h
//  BertyBridgeDemo
//
//  Created by Rémi BARBERO on 29/04/2021.
//

#import <Foundation/Foundation.h>

#import "BertyDevice_darwin.h"

NS_ASSUME_NONNULL_BEGIN

@interface ConnectedPeer : NSObject

@property (nonatomic, assign, nullable) BertyDevice *client;
@property (nonatomic, assign, nullable) BertyDevice *server;
@property (readwrite, getter=isConnected) BOOL connected;

- (BOOL)isClientReady;
- (BOOL)isServerReady;

@end

NS_ASSUME_NONNULL_END
