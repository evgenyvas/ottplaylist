ottplaylist
===========

Generate playlist for ottplayer

Available handlers
------------------

- pomoyka.xspf - playlist from pomoyka.win in xspf format
- acesearch - channels form search.acestream.net in json format

Output formats
--------------

- ott.m3u - compatible with ottplayer

Configuration parameters
------------------------

- port - for server
- types - array with playlist types. Example:
    ```
        "pomoyka.allfon.proxy": {
            "link": "http://pomoyka.win/trash/ttv-list/allfon.all.proxy.xspf",
            "handler": "pomoyka.xspf"
        },
    ```
    - key is type name
    - link - from where download source data
    - handler - see `Available handlers`
- playlists - array of output playlists. Example:
    ```
        {
            "type": "pomoyka.allfon.iproxy",
            "name": "pomoyka.allfon.ott.m3u8",
            "IP": "192.168.1.141",
            "port": 6878,
            "format": "ott.m3u"
        },
    ```
    - type - name from `types` configuration parameter
    - name - route
    - IP - for channels inside playlist
    - port - for channels inside playlist
    - format - see `Output formats`
