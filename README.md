# Trainberry - Puce embarquée

**Ce code est conçu pour le RPi Pico W. Il pourrait ne pas fonctionner sur d'autres cartes.**

Ce projet contient le code à déployer sur les cartes embarquées.

## Quick start

### Configuration

Tout d'abord, vous devez configurer votre carte : IP du train, nom, modèle, paramètres Wi-Fi, etc. Le moyen le plus simple est de copier le fichier internal/configuration/configuration.go.template vers internal/configuration/configuration.go, puis de modifier les valeurs avec les vôtres.

Comme il y a beaucoup de problèmes avec les requêtes ARP sur le Pico W, vous devez également fournir vous-même l'adresse MAC du serveur central, au format hexadécimal. Par exemple, si l'adresse MAC du serveur central est ab:cd:ef:01:02:03, la variable ServerMAC devra être [6]byte{0xAB, 0xCD, 0xEF, 0x01, 0x02, 0x03}.
### Flash

Pour flasher, vous devez spécifier votre cible et une stack-size spécifique :

```sh
tinygo flash -target=pico -stack-size=16kb -monitor  ./cmd
```

Testé avec TinyGo v0.33.0.

## Utilisation générale

Au démarrage, la puce se connecte au réseau Wi-Fi et demande sa propre IP. Pour des raisons de performance, il n’y a pas de traitement réseau "complexe" (pas de DHCP, pas de DNS, etc.).

Une fois connectée, la puce effectue une requête `POST /register` vers le serveur central, avec une charge utile au format JSON ressemblant à ceci :

```json
{
    "name": "BB27000",
    "model": "BB27000",
    "ip": "192.168.1.150"
}
```

La puce est alors prête à accepter de nouvelles connexions pour vérifier son état ou modifier son statut.

## API

La puce expose une API sur le port `80/tcp`. Pour des raisons de performance, il n’y a ni authentification ni autre mécanisme de sécurité. Il est fortement recommandé d’utiliser un réseau Wi-Fi dédié pour votre infrastructure ferroviaire (un vieux routeur Wi-Fi 100 Mbps suffit largement).

L’ensemble de la documentation de référence est disponible dans le dossier `api`.