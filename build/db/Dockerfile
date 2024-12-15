FROM postgres:17.1


ENV DICT="/usr/share/hunspell/ru_RU"
WORKDIR /usr/share/postgresql/17.1/tsearch_data
RUN  apt-get update \
&& apt-get install -y myspell-ru \
&& iconv -f koi8-r -t utf-8 -o russian.affix ${DICT}.aff \
&& iconv -f koi8-r -t utf-8 -o russian.dict  ${DICT}.dic

EXPOSE 5432

