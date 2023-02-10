INSERT INTO product_tags (id, name)
VALUES (1, 'child'), (2, 'adult'), (3, 'luxury'), (4, 'trending'), (5, '18+'), (6, 'expensive');

INSERT INTO countries (code, name, continent_name)
VALUES (84, 'Viet Name', 'Asia'), (86, 'China', 'Asia'), (4, 'America', 'EU');

INSERT INTO users (id, full_name, email, phone, hashed_password)
VALUES (1, 'Vo Nguyen Giap', 'paig@gmail.com', '090909090', '$2a$10$xxlmP1SNdlQyfhKxcam3K.ti6ct9BlkXeL52bF663UAdU3AeAh9Om'),
(2, 'VO Tran kim Chi', 'kchi@gmail.com', '09090909', '$2a$10$LT1eMWPJHBMIiUwfH8ptauviZ2MfrQea3QXcUva88fHN1/wdTCuOS');

INSERT INTO merchants (id, country_code, merchant_name, user_id, description)
VALUES (1, 84, 'Giap Paig', 1, 'none'),
        (2, 86, 'KChi', 2, 'none');

INSERT INTO products (id, name, merchant_id, status)
VALUES (1, 'Ca Phe Rang Bo', 1, 'in_stock'),
        (2, 'Bee Doves', 2, 'running_low'),
        (3, 'Cà Phê Chồn', 1, 'running_low'),
        (4, 'Cà Phê Highland', 1, 'in_stock'),
        (5, 'Pony Pochi', 2, 'out_of_stock');

INSERT INTO product_tags_products (product_tags_id, products_id)
VALUES (6, 4), (6, 5), (6, 3), (2, 3), (1, 5), (4, 1), (4, 2), (4, 3), (4, 4);

INSERT INTO product_pricing (id, product_id, base_price, start_date, end_date)
VALUES (1, 1, 100, '2023-01-10', '2023-02-01'), 
        (2, 2, 70, '2022-01-10', '2023-01-01'),
        (3, 2, 60, '2023-01-10', '2023-02-01'),
        (4, 3, 220, '2023-01-10', '2023-01-12'),
        (5, 4, 160, '2023-01-10', '2023-02-13'),
        (6, 5, 120, '2023-01-10', '2023-02-01');

INSERT INTO deals (id, name, start_date, end_date, type, discount_rate, merchant_id, deal_limit)
VALUES (1, 'Cafe Ngon', '2023-01-10', '2023-01-12', 'discount', 0.01, 1, 25),
        (2, 'Cafe Ngon X1', '2023-02-10', '2023-02-12', 'discount', 0.02, 1, 15),
        (3, 'Cafe Ngon X2', '2022-01-10', '2023-05-12', 'discount', 0.005, 1, 15),
        (4, 'Sale Souvernear', '2023-01-10', '2023-02-12', 'discount', 0.13, 2, 30),
        (5, 'Sale Souvernear X2', '2023-02-10', '2023-02-12', 'discount', 0.22, 2, 40);

INSERT INTO product_general_criteria (id, criteria)
VALUES (1, 'cheap'), (2, 'normal');

INSERT INTO product_size (id, size_value) 
VALUES (1, 'none'), (2, 'small'), (3, 'big');

INSERT INTO product_colour (id, colour_name)
VALUES (1, 'none'), (2, 'pink'), (3, 'blue'), (4, 'yellow');

INSERT INTO product_entry (id, product_id, general_criteria_id, quantity, deal_id, colour_id, size_id)
VALUES (1, 1, 1, 20, 1, 1, 1), (2, 4, 2, 100, 3, 1, 1), (3, 2, 1, 100, 4, 4, 2), (4, 3, 2, 50, 2, 1, 1);
