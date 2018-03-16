import sys
import json
from kapacitor.udf.agent import Agent, Handler, Server
from kapacitor.udf import udf_pb2
import signal

import logging
logging.basicConfig(level=logging.DEBUG, format='%(asctime)s %(levelname)s:%(name)s: %(message)s')
logger = logging.getLogger()


class AgentHandler(Handler):
    def __init__(self, agent):
        self._agent = agent


    def info(self):
        response = udf_pb2.Response()
        response.info.wants = udf_pb2.BATCH
        response.info.provides = udf_pb2.BATCH
        return response

    def init(self, init_req):
        response = udf_pb2.Response()
        response.init.success = True
        return response

    def snapshot(self):
        response = udf_pb2.Response()
        response.snapshot.snapshot = ''
        return response

    def restore(self, restore_req):
        response = udf_pb2.Response()
        response.restore.success = False
        response.restore.error = 'not implemented'
        return response

    def begin_batch(self, begin_req):
        logger.info("Begin Batch")
        response = udf_pb2.Response()
        response.begin.CopyFrom(begin_req)
        self._agent.write_response(response, flush=True)

    def point(self, point):
        response = udf_pb2.Response()
        response.point.CopyFrom(point)
        print "Point.time : %s" %  point.time
        print "Point.tags : %s" %  point.tags["port"]
        print "Point.octets_in : %s" %  point.fieldsDouble["octets_in"]
        print "Point.octets_out : %s" %  point.fieldsDouble["octets_out"]
        print "Point.packets_in : %s" %  point.fieldsDouble["packets_in"]
        print "Point.packets_out : %s" %  point.fieldsDouble["packets_out"]
        self._agent.write_response(response, flush=True)

    def end_batch(self, end_req):
        logger.info("End Batch")
        response = udf_pb2.Response()
        response.end.CopyFrom(end_req)
        self._agent.write_response(response, flush=True)

class accepter(object):
    _count = 0
    def accept(self, conn, addr):
        self._count += 1
        a = Agent(conn, conn)
        h = AgentHandler(a)
        a.handler = h

        logger.info("Starting Agent for connection %d", self._count)
        a.start()
        a.wait()
        logger.info("Agent finished connection %d",self._count)

if __name__ == '__main__':
    path = "/var/lib/kapacitor/agent.sock"
    if len(sys.argv) == 2:
        path = sys.argv[1]
    server = Server(path, accepter())
    logger.info("Started server")
    server.serve()
