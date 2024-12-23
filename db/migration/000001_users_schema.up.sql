create table if not exists tbl_user
(
    id                                   bigserial
        constraint idx_41395_primary
            primary key,
    username                             varchar(50)                                    not null
        constraint ukk0bty7tbcye41jpxam88q5kj2
            unique,
    password                             char(32)                                       not null,
    access_token                         varchar(50)      default NULL::character varying,
    name                                 varchar(128)     default NULL::character varying,
    altname                              varchar(128)     default NULL::character varying,
    address                              varchar(256)     default NULL::character varying,
    email                                varchar(128)     default NULL::character varying,
    status                               smallint                                       not null,
    keyword                              varchar(128)     default ''::character varying,
    shortcode                            bigint,
    description                          varchar(1000)    default NULL::character varying,
    login_attempts                       bigint           default '0'::bigint           not null,
    last_atmpt_time                      timestamp,
    pw_changed_time                      timestamp,
    mobile_no                            varchar(45)      default NULL::character varying,
    acl_list                             jsonb,
    last_login_time                      timestamp,
    create_time                          timestamp,
    update_time                          timestamp,
    push_rate                            bigint           default '10'::bigint,
    mt_port                              varchar(15050)                                 not null,
    opt_out_tag                          varchar(160)     default NULL::character varying,
    vat_registered_no                    varchar(100)     default ''::character varying,
    role_id                              bigint
        constraint fkuser2role
            references tbl_role
            on update cascade on delete restrict,
    approver_id                          bigint
        constraint fkuser2role4approver
            references tbl_role
            on update cascade on delete set null,
    ad_policy                            smallint,
    created_by                           bigint
        constraint fkuser2user4createdby
            references tbl_user
            on update cascade on delete set null,
    accepted_tc                          smallint         default '0'::smallint         not null,
    locked_time                          timestamp,
    old_passwords                        varchar(1000)    default NULL::character varying,
    force_pw_change                      smallint,
    allow_international                  smallint                                       not null,
    rate_card_id                         bigint
        constraint fk3jwmjle5iio57csckuyts0x8v
            references tbl_rate_card,
    user_type                            smallint         default '0'::smallint         not null,
    dealer_id                            bigint,
    admin_id                             bigint,
    b_lead_type                          smallint,
    b_lead_name                          varchar(256)     default NULL::character varying,
    inbox_url                            varchar(255)     default NULL::character varying,
    delivery_url                         varchar(255)     default NULL::character varying,
    inbox_url_type                       varchar(10)      default NULL::character varying,
    delivery_url_type                    varchar(10)      default NULL::character varying,
    obd_port                             varchar(255)     default NULL::character varying,
    camp_id                              bigint,
    masking_enable_for_other_op          smallint         default '0'::smallint         not null,
    industry_type                        smallint,
    lead_id                              varchar(100)     default NULL::character varying,
    pc_rate                              double precision default '1'::double precision not null,
    financial_type                       smallint         default '1'::smallint,
    allow_credit_reimbursement           smallint         default '0'::smallint,
    mid_expiry_time                      double precision default '1'::double precision,
    otp                                  varchar(10)      default NULL::character varying,
    otp_expire_time                      timestamp,
    nid                                  varchar(20)      default NULL::character varying,
    otp_attempts                         bigint           default '0'::bigint,
    vmsisdn                              varchar(20)      default NULL::character varying,
    b_representative_no                  varchar(50),
    b_tin_no                             varchar(100),
    bnid                                 varchar(255),
    business_type                        varchar(255),
    connectivity_ip                      varchar(255),
    connectivity_type                    varchar(255),
    default_long_code                    varchar(255),
    failed_attempt                       integer          default 0,
    force_pw_changed                     integer,
    account_non_locked                   boolean          default true,
    is_client_pending_for_admin_approval boolean,
    is_policy_checked                    boolean,
    kam_one                              varchar(255),
    kam_two                              varchar(255),
    kam_type_one                         varchar(255),
    kam_type_two                         varchar(255),
    last_5_pswd                          varchar(500),
    last_atmt_time                       timestamp,
    mo_port                              varchar(100),
    otp_exp_time_ms                      bigint,
    updated_by                           bigint,
    fin_team_email                       varchar(255),
    data_camp_id                         bigint
);

comment on column tbl_user.password is 'Encrypted password';

alter table tbl_user
    owner to user_adareach_prod;

create index admin_id
    on tbl_user (admin_id);

create index b_lead_type
    on tbl_user (b_lead_type);

create index dealer_id
    on tbl_user (dealer_id);

create index fkuser2role
    on tbl_user (role_id);

create index fkuser2role4approver
    on tbl_user (approver_id);

create index fkuser2user4createdby
    on tbl_user (created_by);

create index user_type
    on tbl_user (user_type);

grant select on tbl_user to user_operation_prod;

