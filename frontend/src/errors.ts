import { currentLang } from './i18n'

type Labels = Record<string, string>

const FIELD_LABELS: Record<'ru' | 'be', Labels> = {
  ru: {
    UserType: 'Тип пользователя',
    NameOrganization: 'Наименование организации',
    Login: 'Логин',
    Password: 'Пароль',
    PasswordConfirm: 'Подтверждение пароля',
    Email: 'Email',
    Phone: 'Телефон',
    LastName: 'Фамилия',
    FirstName: 'Имя',
    Patronymic: 'Отчество',
    IsAbroad: 'Признак «за границей»',
    Locality: 'Населённый пункт',
    Street: 'Улица',
    House: 'Дом',
    Apartment: 'Квартира',
    NoteAddress: 'Примечание к адресу',
    PostalCode: 'Почтовый индекс',
    ApplicantCategoryCode: 'Категория заявителя',
    TopicCode: 'Тема обращения',
    OrganizationName: 'Наименование организации',
    AppealKindCode: 'Вид обращения',
    AppealTypeCode: 'Тип обращения',
    IdentificationNumber: 'Идентификационный номер',
    PhysicalStatusCode: 'Физический статус',
    SocialCategoryCode: 'Социальная категория',
    CitizenshipCode: 'Гражданство',
    Region: 'Регион',
    District: 'Район',
    Content: 'Содержание',
    PlannedRegistrationDate: 'Планируемая дата регистрации',
    StatusCode: 'Статус',
    Language: 'Язык интерфейса',
  },
  be: {
    UserType: 'Тып карыстальніка',
    NameOrganization: 'Найменне арганізацыі',
    Login: 'Лагін',
    Password: 'Пароль',
    PasswordConfirm: 'Пацвярджэнне пароля',
    Email: 'Email',
    Phone: 'Тэлефон',
    LastName: 'Прозвішча',
    FirstName: 'Імя',
    Patronymic: 'Імя па бацьку',
    IsAbroad: 'Прыкмета «за мяжой»',
    Locality: 'Населены пункт',
    Street: 'Вуліца',
    House: 'Дом',
    Apartment: 'Кватэра',
    NoteAddress: 'Заўвага да адрасу',
    PostalCode: 'Паштовы індэкс',
    ApplicantCategoryCode: 'Катэгорыя заяўніка',
    TopicCode: 'Тэма звароту',
    OrganizationName: 'Найменне арганізацыі',
    AppealKindCode: 'Від звароту',
    AppealTypeCode: 'Тып звароту',
    IdentificationNumber: 'Ідэнтыфікацыйны нумар',
    PhysicalStatusCode: 'Фізічны статус',
    SocialCategoryCode: 'Сацыяльная катэгорыя',
    CitizenshipCode: 'Грамадзянства',
    Region: 'Вобласць',
    District: 'Раён',
    Content: 'Змест',
    PlannedRegistrationDate: 'Плануемая дата рэгістрацыі',
    StatusCode: 'Статус',
    Language: 'Мова інтэрфейсу',
  },
}

const TAG_MESSAGES: Record<string, Labels> = {
  required: { ru: 'Поле обязательное', be: 'Поле абавязковае' },
  required_if: {
    ru: 'Обязательное поле для выбранного типа',
    be: 'Абавязковае поле для абранага тыпу',
  },
  email: { ru: 'Некорректный email', be: 'Некарэктны email' },
  phone: { ru: 'Формат +375XXXXXXXXX', be: 'Фармат +375XXXXXXXXX' },
  password: {
    ru: 'Нужны заглавная, строчная буква и цифра',
    be: 'Патрэбны вялікая, малая літара і лічба',
  },
  gte: { ru: 'Слишком короткое значение', be: 'Занадта кароткае значэнне' },
  lte: { ru: 'Слишком длинное значение', be: 'Занадта доўгае значэнне' },
  eqfield: { ru: 'Поля не совпадают', be: 'Палі не супадаюць' },
  oneof: { ru: 'Недопустимое значение', be: 'Недапушчальнае значэнне' },
}

function humanizeMessage(tag: string, lang: 'ru' | 'be'): string {
  const entry = TAG_MESSAGES[tag]
  return entry ? entry[lang] : tag
}

export function humanizeFields(
  fields: Record<string, string> | undefined,
): Record<string, string> {
  if (!fields) return {}
  const lang = currentLang()
  const labels = FIELD_LABELS[lang]
  const out: Record<string, string> = {}
  for (const [key, value] of Object.entries(fields)) {
    const label = labels[key] ?? key
    const hasCyrillic = /[а-яёіў’]/i.test(value)
    out[label] = hasCyrillic ? value : humanizeMessage(value, lang)
  }
  return out
}
