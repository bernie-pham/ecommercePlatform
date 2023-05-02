import pandas as pd
import sqlalchemy
import argparse
import logging
import random


db_driver = 'postgresql+psycopg2://root:secret@localhost:5432/ecommerce_platform?sslmode=disable'
table = 'product_tags'

logging.basicConfig()
logging.getLogger('sqlalchemy.engine').setLevel(logging.INFO)

merchant_mapping = {
    'Arts, Crafts & Sewing': 1,
    'Automotive': 2,
    'Baby Products': 3,
    'Beauty & Personal Care': 6,
    'Clothing, Shoes & Jewelry': 7,
    'Electronics': 4,
    'Grocery & Gourmet Food': 5,
    'Health & Household': 8,
    'Hobbies': 9,
    'Home & Kitchen': 10,
    'Industrial & Scientific': 11,
    'Musical Instruments': 12,
    'Office Products': 13,
    'Patio, Lawn & Garden': 14,
    'Pet Supplies': 15,
    'Remote & App Controlled Vehicles & Parts': 16,
    'Sports & Outdoors': 17,
    'Tools & Home Improvement': 18,
    'Toys & Games': 19,
    'Movies & TV': 1,
    'Cell Phones & Accessories': 4,
    'Video Games': 19,
    'Remote & App Controlled Vehicle Parts': 18
}

colors = [1,2,3,4,5]
sizes = [1,2,3,4,5,6]
quantities = [100, 200, 182, 82, 82, 1090, 20, 40, 10, 5, 70, 1002, 101, 105, 1234, 313, 423, 5453, 12, 432, 54, 27, 84, 11]



def initConn():
    connector = sqlalchemy.create_engine(db_driver)
    return connector


def import_price(dataframe, connector):
    with connector.connect() as conn:
        for idx, value in dataframe.iterrows():
            priceStr = value.values[2].strip()
            if priceStr.startswith('$'):
                if '-' in priceStr:
                    pricevalue = priceStr.split('-')[0].replace(",", "").split('$')[1]
                else:
                    pricevalue = priceStr.replace(",", "").split('$')[1]
            else:
                pricevalue = '$100'.split('$')[1]
            pricevalue = pricevalue.replace(" ", "").strip()
            price = float(pricevalue)
            start_date = '2022-04-11 23:44:43.412 +0700'
            end_date = '2025-04-11 23:44:43.412 +0700'
            priority=1
            sql = sqlalchemy.text(
                "insert into product_pricing \
                (product_id, base_price, start_date, end_date, priority) \
                    values ('{pid}', '{price}', '{sd}', '{ed}', '{prio}');"
                .format(pid=idx, price=price, sd=start_date, ed=end_date, prio=priority))
            conn.execute(sql)
        conn.commit()


def import_tags(dataframe, connector):
    # pre-process
    # de-dup dataframe
    tags = dataframe.drop_duplicates()
    tag_list = {}
    with connector.connect() as conn:
        for idx, tag in tags.iterrows():
            tag_list[tag.values[0].replace("'", "")] = idx
        #     sql = sqlalchemy.text(
        #         """insert into product_tags (id, name) values ('{id}', '{tag}');""".format(tag=tag.values[0].replace("'", ""), id=idx))
        #     conn.execute(sql)
        # conn.commit()
    return tag_list


def import_pro_tags(dataframe, connector, tag_list):
    with connector.connect() as conn:
        for idx, value in dataframe.iterrows():
            tag = value.values[1].replace("'", "")
            tag_id = tag_list[tag]
            sql = sqlalchemy.text(
                "insert into product_tags_products (product_tags_id, products_id) values ('{pid}', '{tid}');"
                .format(pid=tag_id, tid=idx))
            conn.execute(sql)
        conn.commit()


def import_products(dataframe, connector):
    with connector.connect() as conn:
        for idx, value in dataframe.iterrows():
            merchant_id = merchant_mapping[value.values[1].split('|')[0].strip()]
            sql = sqlalchemy.text(
                "insert into products (id, name, status, merchant_id, img_path) values ('{id}', '{name}', '{status}', '{mid}', '{img}');"
                .format(id=idx, name=value.values[0].replace("'", "").replace(".", ""), status='in_stock', mid=merchant_id, img=value.values[3].split('|')[0]))
            conn.execute(sql)
        conn.commit()




def import_pro_entry(rs, connector):
    with connector.connect() as conn:
        for row in rs:
            num_entries = random.randint(1, 3)
            for i in range(num_entries):
                sid = sizes[random.randint(0, 5)]
                cid = colors[random.randint(0, 4)]
                quantity = quantities[random.randint(0, 23)]
                sql = sqlalchemy.text(
                    "insert into product_entry (product_id, colour_id, size_id, quantity) values ('{pid}', '{cid}', '{sid}', '{quantity}');"
                    .format(pid=row.id, cid=cid, sid=sid, quantity=quantity))
                conn.execute(sql)
            conn.commit()



def data_orchestrator(input_file, sheet_name, num_row):
    connector = initConn()
    # First reading all rows from source file, and dedup data, and then import them into table product_tags
    # df_tags = read_data(input_file=input_file,
    #                     sheet_name=sheet_name, chunk_size=None, usecol="A, E")
    # tag_list = import_tags(df_tags, connector)
    # import products
    # last_cursor = 1
    # step = 10
    # while last_cursor < int(num_row):
        # df_products = read_data(input_file=input_file,
        #                         sheet_name=sheet_name, chunk_size=step, usecol="A, B", old_end=last_cursor, idx_col=0)
        # import_products(df_products, connector)
        # import_pro_tags(df_products, connector, tag_list)
        # import_price(df_products, connector)
        # import_pro_entry(df_products, connector)
        # last_cursor += step
    with connector.connect() as conn:
        sql = sqlalchemy.text("select id from products")
        rs = conn.execute(sql).all()
    import_pro_entry(rs, connector)
    




def read_data(input_file, sheet_name, chunk_size=10, old_start=0, old_end=1, usecol="A", idx_col=0):
    try:
        df = pd.read_excel(input_file, sheet_name=sheet_name,
                           nrows=chunk_size, usecols=usecol, skiprows=range(old_start, old_end), index_col=idx_col)
        # df.columns = [c.lower() for c in df.columns]
        df = df.dropna()
    except:
        raise Exception('Cannot read ${file}, with sheetname=${sn}'.format(
            file=input_file, sn=sheet_name))
    return df


def main():
    parser = argparse.ArgumentParser("csv2pgl")
    parser.add_argument("--input", dest="input_file",
                        type=str, help='Enter excel source file', required=True)
    parser.add_argument("--sheet_name", dest="sheet_name",
                        type=str, help='Enter sheet name to process', required=True)
    parser.add_argument("--nrow", dest="num_row",
                        type=str, help='Enter number of row need to process', required=True)
    args = parser.parse_args()

    data_orchestrator(input_file=args.input_file,
                      sheet_name=args.sheet_name, num_row=args.num_row)


def str2bigint(str):
    poststr = ''
    for char in str:
        if not (ord(char) <= 90 and ord(char) >= 65 or ord(char) <= 122 and ord(char) >= 97):
            poststr += char
    return int(poststr)


if __name__ == "__main__":
    main()
