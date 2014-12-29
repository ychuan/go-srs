/*
The MIT License (MIT)

Copyright (c) 2013-2014 winlin

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package server

import (
    "github.com/winlinvip/go-srs/core"
    "github.com/winlinvip/go-srs/rtmp"
    "errors"
)

var FinalStage = errors.New("rtmp final stage")

type commonStage struct {
    logger core.Logger
    conn *rtmp.Conn
}

type connectStage struct {
    commonStage
}

func (cs *connectStage) ConsumeMessage(msg *rtmp.RtmpMessage) (err error) {
    // always expect the connect app message.
    if !msg.Header.IsAmf0Command() && !msg.Header.IsAmf3Command() {
        return
    }

    var pkt rtmp.RtmpPacket
    if pkt,err = cs.conn.Protocol.DecodeMessage(msg); err != nil {
        return
    }

    // got connect app packet
    if pkt,ok := pkt.(*rtmp.RtmpConnectAppPacket); ok {
        if err = cs.conn.Request.Parse(pkt.CommandObject, pkt.Arguments, cs.logger); err != nil {
            cs.logger.Error("parse request from connect app packet failed.")
            return
        }
        cs.logger.Info("rtmp connect app success")

        // discovery vhost, resolve the vhost from config
        // TODO: FIXME: implements it

        // use next stage.
        cs.conn.Stage = NewFinalStage(cs.conn)
    }

    return
}

type finalStage struct {
    commonStage
}

func NewFinalStage(conn *rtmp.Conn) rtmp.Stage {
    return &finalStage{
        commonStage:commonStage{
            logger:conn.Logger,
            conn: conn,
        },
    }
}

func (fs *finalStage) ConsumeMessage(msg *rtmp.RtmpMessage) (err error) {
    return FinalStage
}