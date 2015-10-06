from flask import Flask
from flask import request
from flask import render_template

import redis


app = Flask(__name__)
r = redis.StrictRedis(host='api.ticketea.dev', port=6379, db=0)


@app.route('/', methods=['GET'])
def index():
  return render_template('panel.html')


@app.route('/update_ticket', methods=['POST'])
def update_sold_tickets():
  fake_begin_key = 'Organizer:1234:Event:3456:Session:246'

  ticket_id = request.form['ticket_id']
  price = float(request.form['price'])

  pipe = r.pipeline()

  # Add ticket to total quantity
  tickets_quantity_key = '%(fake_begin_key)s:TicketsQuantity' % locals()
  pipe.incr(tickets_quantity_key)

  # Add ticket to ticket type quantity
  ticket_type_quantity_key = '%(fake_begin_key)s:TicketType:%(ticket_id)s:quantity' % locals()
  pipe.incr(ticket_type_quantity_key)

  # Add price to ticket type total amount
  ticket_type_amount_key = '%(fake_begin_key)s:TicketType:%(ticket_id)s:amount' % locals()
  pipe.incrbyfloat(ticket_type_amount_key, price)

  pipe.execute()


  # Publish to the organizer channel
  publish_channel = "%(fake_begin_key)s:UpdateCounters" % locals()
  r.publish(publish_channel, True)

  return "OK"




if __name__ == '__main__':
  # 0.0.0.0 para aceptar conexiones externas
  # debug=True para hacer autoreload cuando se guarda
  app.run(host='0.0.0.0', debug=True)