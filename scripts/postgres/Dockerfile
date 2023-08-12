FROM postgres:15

# Copy the database initialization script:
COPY ./init.sh /docker-entrypoint-initdb.d/


RUN apt-get update && apt-get install hunspell-pt-br -y

# RUN iconv -f latin1 -t UTF-8 /usr/share/hunspell/pt_BR.aff | tee /usr/share/postgresql/15/tsearch_data/brazilian.affix
# RUN iconv -f latin1 -t UTF-8 /usr/share/hunspell/pt_BR.dic | tee /usr/share/postgresql/15/tsearch_data/brazilian.dict

RUN tee /usr/share/postgresql/15/tsearch_data/brazilian.affix
RUN tee /usr/share/postgresql/15/tsearch_data/brazilian.dict
