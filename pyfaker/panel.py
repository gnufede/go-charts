from flask import Flask
from flask import request
from flask import render_template
from datetime import date, datetime

import redis


app = Flask(__name__)
r = redis.StrictRedis(host='localhost', port=6379, db=0)


@app.route('/', methods=['GET'])
def index():
  return render_template('panel.html')


@app.route('/update_ticket', methods=['POST'])
def update_sold_tickets():
  fake_begin_key = 'Organizer:1:Event:1:Channel:1:Session:1' #:Zone:1'

  ticket_id = request.form['ticket_id']
  today = date.today().strftime('%Y-%m-%d')
  price = request.form['price']

  pipe = r.pipeline()

  # Add ticket to total quantity
  #tickets_quantity_key = '%(fake_begin_key)s:%(today)s:TicketsQuantity' % locals()
  tickets_quantity_key = '%(fake_begin_key)s:Date:%(today)s' % locals()
  pipe.hincrby(tickets_quantity_key, datetime.now().strftime('%H:%M'), 1)

  # Add ticket to ticket type quantity
  #ticket_type_quantity_key = '%(fake_begin_key)s:TicketType:%(today)s:%(ticket_id)s:quantity' % locals()
  ticket_type_quantity_key = '%(fake_begin_key)s:TicketType:%(ticket_id)s:Date:%(today)s' % locals()
  pipe.hincrby(ticket_type_quantity_key, datetime.now().strftime('%H:%M'), 1)

  # Add price to ticket type total amount
  ticket_type_amount_key = '%(fake_begin_key)s:TicketType:%(ticket_id)s:Date:%(today)s:Amount' % locals()
  pipe.hincrby(ticket_type_amount_key, datetime.now().strftime('%H:%M'), price)

  pipe.execute()


  # Publish to the organizer channel
  publish_channel = "1" #"%(fake_begin_key)s:UpdateCounters" % locals()
  r.publish(publish_channel, True)

  return "OK"




if __name__ == '__main__':
  # 0.0.0.0 para aceptar conexiones externas
  # debug=True para hacer autoreload cuando se guarda
  app.run(host='0.0.0.0', debug=True)
