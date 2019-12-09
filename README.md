# Presidenten

Presidenten is een command-line versie van [het gelijknamige kaartspel](https://nl.wikipedia.org/wiki/Presidenten_(kaartspel)). Het laat spelers toe om met elkaar te spelen via hun eigen computer. Eén van de spelers is de "server" (of "host") en de anderen zijn de "clients" (of "guests") die met hem verbinden.

De applicatie bevat ook de mogelijkheid om een computerprogramma te laten meespelen, genaamd een "bot". Deze zal zelf beslissingen nemen en het spel gewoon meespelen alsof het een menselijke speler is. Het is ook mogelijk om uitsluitend bots tegen elkaar te laten spelen, al dan niet met een zekere vertraging zodat mensen kunnen meekijken.

## Gebruik

    presidenten commando [argumenten]

### Beschikbare combinaties

    presidenten

Start het spel (als mens) door te verbinden met een server.
Geef na het starten het IP-adres op van de server, optioneel
gevolgd door een poortnummer, bv. '192.168.0.100:8080'.

---

    presidenten server [poort]

Start een server voor een presidentenspel. Het poortnummer
kan optioneel worden meegegeven (zo niet wordt 3009 gebruikt).
Om vervolgens te starten met de spelers in de lobby, typ 'start'.
Indien de spelers enkel bots zijn, kan ook 'run' worden getypt,
wat de artificiële vertragingen voor bot-beslissingen uitzet.

---
    presidenten bot <bestandsnaam> [commando]

Start het spel als bot. Dit betekent dat alle beslissingen
voor deze speler worden genomen door een computerprogramma.
Geef de bestandsnaam op van het programma, optioneel gevolgd
door het commando om dat programma uit te voeren. Voor een
Python-programma kan dit bijvoorbeeld 'bot.py python' zijn.
Het computerprogramma zal communiceren met het spel via zijn
standaardinvoer en standaarduitvoer, zie `presidenten pybot`.
Het IP-adres om te verbinden met een server en de spelersnaam
moeten bij de start wel eerst nog manueel worden ingegeven.

---
    presidenten pybot [bestandsnaam]

Genereert een voorbeeld van een computerprogramma dat kan
worden gebruikt als functionele bot in een presidentenspel.
Om dit programma te kunnen gebruiken is [Python 3](https://www.python.org/downloads/) vereist.
Merk echter op dat een bot in gelijk welke programmeertaal kan
worden geschreven, zolang het protocol wordt geïmplementeerd.

---
    presidenten help

Toont een overzicht van de beschikbare commando's.
