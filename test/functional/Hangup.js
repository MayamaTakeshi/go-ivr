const sip = require ('sip-lab')
const Zester = require('zester')
const m = require('data-matching')
const sip_msg = require('sip-matching')

const http = require('http')

sip.set_log_level(6)

const z = new Zester()

async function test() {
	const flags = 0

	sip.dtmf_aggregation_on(500)

	z.trap_events(sip.event_source, 'event', (evt) => {
		var e = evt.args[0]
		return e
	})

	console.log(sip.start((data) => { console.log(data)} ))

	const server = http.createServer(z.callback_trap('http_request', function preprocessor(evt) {
		console.log("evt:")
		console.dir(evt)
		return {
			name: 'http_request',
			req: evt.args[0],
			res: evt.args[1],
		}
	}))

	server.listen(80)

	var t1 = sip.transport_create("127.0.0.1", 5090, 1)

	console.log("t1", t1)

	z.add_event_filter({
		event: 'response',
		method: 'INVITE',
		msg: sip_msg({
			'$rs': '100',
		}),
	})

	var oc = sip.call_create(t1.id, flags, 'sip:0312341234@t', 'sip:0911112222@127.0.0.1:5160')

	await z.wait([
		{    
			name: 'http_request', 
			req: m.collect('req'),
			res: m.collect('res'),
		}
	], 500)

	var res = z.store.res	

	res.writeHead(200, {'Content-Type': 'plain/xml'})
	res.end('<IVR><Hangup cause="USER_BUSY"/></IVR>')

	await z.wait([
		{
			event: 'response',
			call_id: oc.id,
			method: 'INVITE',
			msg: sip_msg({
				$rs: '486',
				$rr: 'Busy Here',
			}),
		},
		{
			event: 'call_ended',
			call_id: oc.id,
		},
	], 500)

	sip.stop()
}


test()
.catch(e => {
    console.error(e)
    process.exit(1)
})
