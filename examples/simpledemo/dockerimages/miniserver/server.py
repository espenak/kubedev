import http.server
import socketserver
import sys

PORT = 8000


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        print(f'DEMO: {self.path}', flush=True)
        message = 'DEMO ffs!'
        self.send_response(200)
        self.send_header('Content-Type',
                         'text/plain; charset=utf-8')
        self.end_headers()
        self.wfile.write(message.encode('utf-8'))
        sys.stderr.flush()


with socketserver.TCPServer(("", PORT), Handler) as httpd:
    print("serving at port", PORT, flush=True)
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print('Bye!', flush=True)
