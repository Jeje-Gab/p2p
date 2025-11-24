-- Update skin image URLs to use MinIO storage
-- This migration updates existing image URLs to point to the MinIO bucket

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/AK-47/Redline.svg',
    updated_at = NOW()
WHERE name = 'Redline' AND weapon_type = 'AK-47';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/AWP/Asiimov.svg',
    updated_at = NOW()
WHERE name = 'Asiimov' AND weapon_type = 'AWP';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/M4A4/Howl.svg',
    updated_at = NOW()
WHERE name = 'Howl' AND weapon_type = 'M4A4';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/AWP/Dragon-Lore.svg',
    updated_at = NOW()
WHERE name = 'Dragon Lore' AND weapon_type = 'AWP';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/USP-S/Kill-Confirmed.svg',
    updated_at = NOW()
WHERE name = 'Kill Confirmed' AND weapon_type = 'USP-S';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/Desert-Eagle/Printstream.svg',
    updated_at = NOW()
WHERE name = 'Printstream' AND weapon_type = 'Desert Eagle';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/Desert-Eagle/Hypnotic.svg',
    updated_at = NOW()
WHERE name = 'Hypnotic' AND weapon_type = 'Desert Eagle';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/Glock-18/Fade.svg',
    updated_at = NOW()
WHERE name = 'Fade' AND weapon_type = 'Glock-18';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/M4A4/The-Emperor.svg',
    updated_at = NOW()
WHERE name = 'The Emperor' AND weapon_type = 'M4A4';

UPDATE p2p.skins SET
    image_url = 'http://localhost:9000/net.public.p2p/skins/Desert-Eagle/Neo-Noir.svg',
    updated_at = NOW()
WHERE name = 'Neo-Noir' AND weapon_type = 'Desert Eagle';

COMMENT ON TABLE p2p.skins IS 'Skin images are now stored in MinIO bucket: net.public.p2p';
