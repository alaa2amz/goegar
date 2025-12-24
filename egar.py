from datetime import date
from csv import DictReader
from os.path import join as path_join
from os.path import isfile as exist
from pathlib import Path
from pprint import pp

data_file_name = 'egar.csv'
prefix = ''
water_dir = 'water'
elec_dir = 'elec'
fix_dir = 'fix'
message_dir = 'message'

stair_elec_plate = '790'
basement_elec_plate = '789'
stair_clean_fare = 50
water_factor = 6.5
factor = 7.0

elec_dict = {}
today = date.today()
message = ''

#TODO: add message
def format_egar(data_dict,elec_dict,water_cost=0,maintinance=0):
    dd = data_dict
    ed = elec_dict
    name = dd['name']
    mansion = dd['mansion']
    plate = dd['elec']
    rent = int(dd['egar'])
    elec = ed.get(plate, 0)
    water = water_cost / water_factor
    stair = int(ed[stair_elec_plate]) / factor
    base = int(ed[basement_elec_plate]) / factor
    clean = stair_clean_fare
    fix = maintinance // factor
    total = rent + elec + water + stair + base + clean + fix
    temp = f'''
===========================
name:       Mr. {name}
mansion:    -- {mansion}
plate:      -- {plate}
rent:       {rent}
elec:       {elec}
water:      {water:.2f}
stair:      {stair:.2f}
base:       {base:.2f}
clean:      {clean}
service:    {fix:.2f}

--------------------------
total:   {total:.2f}
==============================
'''
    return temp

#TODO: date = date.today()
#TODO: maintainwnce and messages
def initdata(date):
    water_cost = 0
    year = str(date.year)
    month = str(date.month)
    elec_file_path = path_join(
            prefix, elec_dir, year, month)
    water_file_path = path_join(
            prefix, water_dir, year, month)
    #Path(water_file_path).touch()

    if exist(water_file_path): 
        with open(water_file_path) as water_stream:
            water_cost = int(water_stream.readall())

    with open(elec_file_path) as elec_stream:
        for line in elec_stream:
            fields = line.split()
            if len(fields) < 2:
                continue
            elec_dict[fields[0]] = int(fields[1])

    #debug prints
    print(date.month)
    print('water data')
    print(water_cost)
    print('elec data')
    pp(elec_dict)

    
    with open(data_file_name) as data_stream:
        data_dict_list = list(
                DictReader(
                    data_stream, dialect='excel-tab',))

    #debug prints
    print('data')
    pp(data_dict_list)

    result = ''
    for data_dict in data_dict_list:
        result += format_egar(data_dict,elec_dict)
    print(result)
    

if __name__ == '__main__':
    initdata(today)
