require "excon"

resp1 = Excon.get("http://example.com/")
resp1.body

resp2 = Excon.post("http://example.com/", body: "some_data")
resp2.body

resp3 = Excon.put("http://example.com/", body: "some_data")
resp3.body

resp4 = Excon.patch("http://example.com/", body: "some_data")
resp4.body

resp5 = Excon.delete("http://example.com/")
resp5.body

resp6 = Excon.head("http://example.com/")
resp6.body

resp7 = Excon.options("http://example.com/")
resp7.body

resp8 = Excon.trace("http://example.com/")
resp8.body

connection = Excon.new("http://example.com")
resp9 = connection.get(path: "/")
resp9.body

resp10 = connection.post(path: "/foo")
resp10.body